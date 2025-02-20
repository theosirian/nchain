package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	dbconf "github.com/kthomas/go-db-config"
	pgputil "github.com/kthomas/go-pgputil"
	redisutil "github.com/kthomas/go-redisutil"

	"github.com/provideplatform/nchain/common"
	"github.com/provideplatform/nchain/connector"
)

const runloopTickerInterval = 5 * time.Second
const runloopSleepInterval = 250 * time.Millisecond

var (
	cancelF     context.CancelFunc
	closing     uint32
	shutdownCtx context.Context

	mutex sync.Mutex

	connectors []*connector.Connector
)

func init() {
	if common.ConsumeNATSStreamingSubscriptions {
		common.Log.Panicf("reachabilitydaemon instance started with CONSUME_NATS_STREAMING_SUBSCRIPTIONS=true")
		return
	}

	pgputil.RequirePGP()
	redisutil.RequireRedis()
}

func main() {
	common.Log.Debug("Installing signal handlers for reachabilitydaemon")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	shutdownCtx, cancelF = context.WithCancel(context.Background())

	requireConnectorReachabilityDaemonInstances()

	common.Log.Debugf("Running reachabilitydaemon main()")
	timer := time.NewTicker(runloopTickerInterval)
	defer timer.Stop()

	for !shuttingDown() {
		select {
		case <-timer.C:
			// TODO: check reachability and statsdaemon statuses
		case sig := <-sigs:
			common.Log.Infof("Received signal: %s", sig)
			shutdown()
		case <-shutdownCtx.Done():
			close(sigs)
		default:
			time.Sleep(runloopSleepInterval)
		}
	}

	common.Log.Debug("Exiting reachabilitydaemon main()")
	cancelF()
}

func requireConnectorReachabilityDaemonInstances() {
	mutex.Lock()
	defer mutex.Unlock()

	db := dbconf.DatabaseConnection()

	connectors = make([]*connector.Connector, 0)
	db.Find(&connectors)

	for _, connector := range connectors {
		host := connector.Host(db)

		if host != nil {
			cfg := connector.ParseConfig()
			port, portOk := cfg["port"].(float64)
			apiPort, apiPortOk := cfg["api_port"].(float64)

			reachableFn := func() {
				connector.UpdateStatus(db, "available", nil)
			}

			unreachableFn := func() {
				connector.Reload(db)

				if connector.Status != nil && *connector.Status == "deprovisioning" {
					if portOk {
						EvictReachabilityDaemon(&endpoint{
							network: "tcp",
							addr:    fmt.Sprintf("%s:%v", *host, port),
						})
					}
					if apiPortOk {
						EvictReachabilityDaemon(&endpoint{
							network: "tcp",
							addr:    fmt.Sprintf("%s:%v", *host, port),
						})
					}
				} else {
					connector.UpdateStatus(db, "unreachable", nil)
				}
			}

			if portOk {
				RequireReachabilityDaemon(&endpoint{
					network:     "tcp",
					addr:        fmt.Sprintf("%s:%v", *host, port),
					reachable:   reachableFn,
					unreachable: unreachableFn,
				})
			}

			if apiPortOk {
				RequireReachabilityDaemon(&endpoint{
					network:     "tcp",
					addr:        fmt.Sprintf("%s:%v", *host, apiPort),
					reachable:   reachableFn,
					unreachable: unreachableFn,
				})
			}
		}
	}
}

func shutdown() {
	if atomic.AddUint32(&closing, 1) == 1 {
		common.Log.Debug("Shutting down reachabilitydaemon")
		cancelF()
	}
}

func shuttingDown() bool {
	return (atomic.LoadUint32(&closing) > 0)
}
