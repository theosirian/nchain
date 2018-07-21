package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/provideservices/provide-go"

	"github.com/gorilla/websocket"
	logger "github.com/kthomas/go-logger"
)

const defaultStatsDaemonQueueSize = 8
const networkStatsJsonRpcPollingTickerInterval = time.Millisecond * 2500
const networkStatsMaxRecentBlockCacheSize = 32
const networkStatsMinimumRecentBlockCacheSize = 3
const statsDaemonMaximumBackoffMillis = 12800

var currentNetworkStats = map[string]*StatsDaemon{}
var currentNetworkStatsMutex = &sync.Mutex{}

type NetworkStatsDataSource struct {
	Network *Network
	Poll    func(chan *provide.NetworkStatus) error                         // JSON-RPC polling -- implementations should be blocking
	Stream  func(chan *provide.EthereumWebsocketSubscriptionResponse) error // websocket -- implementations should be blocking
}

type StatsDaemon struct {
	attempt    uint32
	backoff    int64
	dataSource *NetworkStatsDataSource

	log *logger.Logger

	cancelF     context.CancelFunc
	closing     uint32
	shutdownCtx context.Context

	queue       chan *provide.EthereumWebsocketSubscriptionResponse
	statusQueue chan *provide.NetworkStatus

	recentBlocks          []interface{}
	recentBlockTimestamps []uint64
	stats                 *provide.NetworkStatus
}

type jsonRpcNotSupported string
type websocketNotSupported string

func (err jsonRpcNotSupported) Error() string {
	return "JSON-RPC not supported"
}

func (err websocketNotSupported) Error() string {
	return "Websocket not supported"
}

// NetworkStatsDataSourceFactory builds and returns a JSON-RPC and streaming websocket
// data source which is used by stats daemon instances to consume network statistics
func NetworkStatsDataSourceFactory(network *Network) *NetworkStatsDataSource {
	return &NetworkStatsDataSource{
		Network: network,

		Poll: func(ch chan *provide.NetworkStatus) error {
			rpcURL := network.rpcURL()
			if rpcURL == "" {
				err := new(jsonRpcNotSupported)
				return *err
			}
			ticker := time.NewTicker(networkStatsJsonRpcPollingTickerInterval)
			for {
				select {
				case <-ticker.C:
					status, err := network.Status(true)
					if err != nil {
						Log.Errorf("Failed to retrieve network status via JSON-RPC: %s; %s", rpcURL, err)
						ticker.Stop()
						return nil
					}

					ch <- status
				}
			}
		},

		Stream: func(ch chan *provide.EthereumWebsocketSubscriptionResponse) error {
			websocketURL := network.websocketURL()
			if websocketURL == "" {
				err := new(websocketNotSupported)
				return *err
			}
			var wsDialer websocket.Dialer
			wsConn, _, err := wsDialer.Dial(websocketURL, nil)
			if err != nil {
				Log.Errorf("Failed to establish network stats websocket connection to %s", websocketURL)
			} else {
				defer wsConn.Close()
				var chainID string
				if network.ChainID != nil {
					chainID = *network.ChainID
				}
				payload := map[string]interface{}{
					"method":  "eth_subscribe",
					"params":  []string{"newHeads"},
					"id":      chainID,
					"jsonrpc": "2.0",
				}
				if err := wsConn.WriteJSON(payload); err != nil {
					Log.Errorf("Failed to write subscribe message to network stats websocket connection")
				} else {
					Log.Debugf("Subscribed to network stats websocket: %s", websocketURL)

					for {
						_, message, err := wsConn.ReadMessage()
						if err != nil {
							Log.Errorf("Failed to receive message on network stats websocket; %s", err)
							break
						} else {
							Log.Debugf("Received message on network stats websocket: %s", message)
							response := &provide.EthereumWebsocketSubscriptionResponse{}
							err := json.Unmarshal(message, response)
							if err != nil {
								Log.Warningf("Failed to unmarshal message received on network stats websocket: %s; %s", message, err.Error())
							} else {
								ch <- response
							}
						}
					}
				}
			}
			return err
		},
	}
}

// Consume the websocket stream; attempts to fallback to JSON-RPC if websocket stream fails or is not available for the network
func (sd *StatsDaemon) consume() []error {
	errs := make([]error, 0)
	Log.Debugf("Attempting to consume configured stats daemon data source; attempt #%v", sd.attempt)
	err := sd.dataSource.Stream(sd.queue)
	if err != nil {
		errs = append(errs, err)
		switch err.(type) {
		case jsonRpcNotSupported:
			sd.log.Warningf("Configured stats daemon data source does not support JSON-RPC: %s; attempting to upgrade to websocket stream for network id: %s", sd.dataSource.Network.ID)
			err := sd.dataSource.Stream(sd.queue)
			if err != nil {
				errs = append(errs, err)
				sd.log.Warningf("Configured stats daemon data source returned error while consuming JSON-RPC endpoint: %s; restarting stream...", err.Error())
			}
		case websocketNotSupported:
			sd.log.Warningf("Configured stats daemon data source does not support streaming via websocket; attempting to fallback to JSON-RPC long polling using stats daemon for network id: %s", sd.dataSource.Network.ID)
			err := sd.dataSource.Poll(sd.statusQueue)
			if err != nil {
				errs = append(errs, err)
				sd.log.Warningf("Configured stats daemon data source returned error while consuming JSON-RPC endpoint: %s; restarting stream...", err.Error())
			}
		}
	}
	return errs
}

func (sd *StatsDaemon) ingest(response interface{}, realtime bool) {
	if sd.dataSource.Network.isEthereumNetwork() {
		switch response.(type) {
		case *provide.EthereumWebsocketSubscriptionResponse:
			resp := response.(*provide.EthereumWebsocketSubscriptionResponse)
			if result, ok := resp.Params["result"].(map[string]interface{}); ok {
				if _, mixHashOk := result["mixHash"]; !mixHashOk {
					result["mixHash"] = common.HexToHash("0x")
				}
				if _, nonceOk := result["nonce"]; !nonceOk {
					result["nonce"] = types.EncodeNonce(0)
				}
				if resultJSON, err := json.Marshal(result); err == nil {
					header := &types.Header{}
					err := json.Unmarshal(resultJSON, header)
					if err != nil {
						Log.Warningf("Failed to stringify result JSON in otherwise valid message received on network stats websocket: %s; %s", response, err.Error())
					} else if header != nil && header.Number != nil {
						sd.ingest(header, realtime)
					}
				}
			}
		case *provide.NetworkStatus:
			resp := response.(*provide.NetworkStatus)
			if resp != nil && resp.Meta != nil {
				if header, headerOk := resp.Meta["last_block_header"].(map[string]interface{}); headerOk {
					if _, mixHashOk := header["mixHash"]; !mixHashOk {
						header["mixHash"] = common.HexToHash("0x")
					}
					if _, nonceOk := header["nonce"]; !nonceOk {
						header["nonce"] = types.EncodeNonce(0)
					}

					if headerJSON, err := json.Marshal(header); err == nil {
						hdr := &types.Header{}
						err := json.Unmarshal(headerJSON, hdr)
						if err != nil {
							Log.Warningf("Failed to stringify result JSON in otherwise valid message received via JSON-RPC: %s; %s", response, err.Error())
						} else if hdr != nil && hdr.Number != nil {
							sd.ingest(hdr, realtime)
						}
					}
				} else {
					Log.Warningf("Failed to parse last_block_header from *provide.NetworkStats meta; dropping message...")
				}
			} else {
				Log.Warningf("Received malformed *provide.NetworkStats message; dropping message...")
			}
		case *types.Header:
			header := response.(*types.Header)
			sd.stats.Block = header.Number.Uint64()
			sd.stats.State = nil
			sd.stats.Syncing = sd.stats.Block == 0

			if sd.stats.Block == 0 {
				Log.Debugf("Ignoring genesis header")
				return
			}

			var lastBlockAt uint64
			if realtime {
				lastBlockAt = uint64(time.Now().UnixNano() / 1000000)
			} else {
				lastBlockAt = header.Time.Uint64() * 1000.0
			}
			sd.stats.LastBlockAt = &lastBlockAt

			sd.stats.Meta["last_block_header"] = header

			if len(sd.recentBlocks) == 0 || sd.recentBlocks[len(sd.recentBlocks)-1].(*types.Header).Hash().String() != header.Hash().String() {
				sd.recentBlocks = append(sd.recentBlocks, header)
				sd.recentBlockTimestamps = append(sd.recentBlockTimestamps, lastBlockAt)
			}

			for len(sd.recentBlocks) > networkStatsMaxRecentBlockCacheSize {
				i := len(sd.recentBlocks) - 1
				sd.recentBlocks = append(sd.recentBlocks[:i], sd.recentBlocks[i+1:]...)
			}

			if len(sd.recentBlocks) >= networkStatsMinimumRecentBlockCacheSize {
				blocktimes := make([]float64, 0)
				timedelta := float64(0)
				i := 0
				for i < len(sd.recentBlocks)-1 {
					currentBlocktime := sd.recentBlockTimestamps[i]
					nextBlocktime := sd.recentBlockTimestamps[i+1]
					blockDelta := float64(nextBlocktime-currentBlocktime) / 1000.0
					blocktimes = append(blocktimes, blockDelta)
					timedelta += blockDelta
					i++
				}

				if len(blocktimes) > 0 {
					sd.stats.Meta["average_blocktime"] = timedelta / float64(len(blocktimes))
					sd.stats.Meta["blocktimes"] = blocktimes
					sd.stats.Meta["last_block_hash"] = header.Hash().String()
				}
			}
		}
	}
}

// This loop is responsible for processing new messages received by daemon
func (sd *StatsDaemon) loop() error {
	for {
		select {
		case msg := <-sd.queue:
			sd.ingest(msg, true)

		case msg := <-sd.statusQueue:
			sd.ingest(msg, false)

		case <-sd.shutdownCtx.Done():
			sd.log.Debugf("Closing stats daemon on shutdown")
			return nil
		}
	}
}

// EvictNetworkStatsDaemon evicts a single, previously-initialized stats daemon instance {
func EvictNetworkStatsDaemon(network *Network) error {
	if daemon, ok := currentNetworkStats[network.ID.String()]; ok {
		Log.Debugf("Evicting stats daemon instance for network: %s; id: %s", *network.Name, network.ID)
		daemon.shutdown()
		currentNetworkStatsMutex.Lock()
		delete(currentNetworkStats, network.ID.String())
		currentNetworkStatsMutex.Unlock()
		return nil
	}
	return fmt.Errorf("Unable to evict stats daemon instance for network: %s; id; %s", *network.Name, network.ID)
}

// RequireNetworkStatsDaemon ensures a single stats daemon instance is running for
// the given network; if no stats daemon instance has been started for the network,
// the instance is configured and started immediately, caching real-time network stats.
func RequireNetworkStatsDaemon(network *Network) *StatsDaemon {
	if network.AvailablePeerCount() == 0 {
		Log.Debugf("Stats daemon instance not initialized for network: %s; no available peers", *network.Name)
		return nil
	}

	var daemon *StatsDaemon
	if daemon, ok := currentNetworkStats[network.ID.String()]; ok {
		Log.Debugf("Cached stats daemon instance found for network: %s; id: %s", *network.Name, network.ID)
		return daemon
	}

	Log.Infof("Initializing new stats daemon instance for network: %s; id: %s", *network.Name, network.ID)
	daemon = NewNetworkStatsDaemon(Log, network)
	currentNetworkStats[network.ID.String()] = daemon

	go daemon.run()

	return daemon
}

// NewNetworkStatsDaemon initializes a new network stats daemon instance using
// NetworkStatsDataSourceFactory to construct the daemon's its data source
func NewNetworkStatsDaemon(lg *logger.Logger, network *Network) *StatsDaemon {
	sd := new(StatsDaemon)
	sd.attempt = 0
	sd.log = lg.Clone()
	sd.shutdownCtx, sd.cancelF = context.WithCancel(context.Background())
	sd.dataSource = NetworkStatsDataSourceFactory(network)
	sd.queue = make(chan *provide.EthereumWebsocketSubscriptionResponse, defaultStatsDaemonQueueSize)
	sd.statusQueue = make(chan *provide.NetworkStatus, defaultStatsDaemonQueueSize)
	//sd.handleSignals()

	chainID := network.ChainID
	if chainID == nil {
		_chainID := hexutil.EncodeBig(provide.GetChainID(network.ID.String(), network.rpcURL()))
		chainID = &_chainID
	}
	sd.stats = &provide.NetworkStatus{
		ChainID: chainID,
		Meta:    map[string]interface{}{},
		State:   stringOrNil("configuring"),
	}

	return sd
}

// Run the configured stats daemon instance
func (sd *StatsDaemon) run() error {
	go func() {
		for !sd.shuttingDown() {
			sd.attempt++
			Log.Debugf("Stepping into main runloop of stats daemon instance; attempt #%v", sd.attempt)
			errs := sd.consume()
			if len(errs) > 0 {
				Log.Warningf("Configured stats daemon data source returned %v error(s) while attempting to consume configured data source", len(errs))
				if sd.backoff == 0 {
					sd.backoff = 100
				} else {
					sd.backoff *= 2
				}
				if sd.backoff > statsDaemonMaximumBackoffMillis {
					sd.backoff = 0
				}
				time.Sleep(time.Duration(sd.backoff) * time.Millisecond)
				sd.dataSource.Network.Reload()
				sd.dataSource.Network.resolveAndBalanceJsonRpcAndWebsocketUrls(DatabaseConnection())
			}
		}
	}()

	err := sd.loop()

	if err == nil {
		sd.log.Info("Stats daemon exited cleanly")
	} else {
		if !sd.shuttingDown() {
			Log.Errorf("Forcing shutdown of stats daemon due to error; %s", err)
			sd.shutdown()
		}
	}

	return err
}

func (sd *StatsDaemon) handleSignals() {
	Log.Debug("Installing SIGINT and SIGTERM signal handlers")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigs:
			Log.Infof("Received signal: %s", sig)
			sd.shutdown()
		case <-sd.shutdownCtx.Done():
			close(sigs)
		}
	}()
}

func (sd *StatsDaemon) shutdown() {
	if atomic.AddUint32(&sd.closing, 1) == 1 {
		Log.Debugf("Shutting down stats daemon instance for network: %s", *sd.dataSource.Network.Name)
		sd.cancelF()
	}
}

func (sd *StatsDaemon) shuttingDown() bool {
	return (atomic.LoadUint32(&sd.closing) > 0)
}
