package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	dbconf "github.com/kthomas/go-db-config"
	natsutil "github.com/kthomas/go-natsutil"
	uuid "github.com/kthomas/go.uuid"
	"github.com/nats-io/go-nats-streaming"
)

const natsTxSubject = "goldmine.tx"
const natsTxMaxInFlight = 2048
const natsTxReceiptSubject = "goldmine.tx.receipt"
const natsTxReceiptMaxInFlight = 2048

const txAckWait = time.Second * 10
const txReceiptAckWait = time.Second * 10

func createNatsTxSubscriptions(natsConnection stan.Conn, wg *sync.WaitGroup) {
	for i := uint64(0); i < natsutil.GetNatsConsumerConcurrency(); i++ {
		wg.Add(1)
		go func() {
			defer natsConnection.Close()

			txSubscription, err := natsConnection.QueueSubscribe(natsTxSubject, natsTxSubject, consumeTxMsg, stan.SetManualAckMode(), stan.AckWait(txAckWait), stan.MaxInflight(natsTxMaxInFlight), stan.DurableName(natsTxSubject))
			if err != nil {
				Log.Warningf("Failed to subscribe to NATS subject: %s", natsTxSubject)
				wg.Done()
				return
			}
			defer txSubscription.Unsubscribe()
			Log.Debugf("Subscribed to NATS subject: %s", natsTxSubject)

			wg.Wait()
		}()
	}
}

func createNatsTxReceiptSubscriptions(natsConnection stan.Conn, wg *sync.WaitGroup) {
	for i := uint64(0); i < natsutil.GetNatsConsumerConcurrency(); i++ {
		wg.Add(1)
		go func() {
			defer natsConnection.Close()

			txReceiptSubscription, err := natsConnection.QueueSubscribe(natsTxReceiptSubject, natsTxReceiptSubject, consumeTxReceiptMsg, stan.SetManualAckMode(), stan.AckWait(txReceiptAckWait), stan.MaxInflight(natsTxReceiptMaxInFlight), stan.DurableName(natsTxReceiptSubject))
			if err != nil {
				Log.Warningf("Failed to subscribe to NATS subject: %s", natsTxReceiptSubject)
				wg.Done()
				return
			}
			defer txReceiptSubscription.Unsubscribe()
			Log.Debugf("Subscribed to NATS subject: %s", natsTxReceiptSubject)

			wg.Wait()
		}()
	}
}

func consumeTxMsg(msg *stan.Msg) {
	Log.Debugf("Consuming %d-byte NATS tx message on subject: %s", msg.Size(), msg.Subject)

	execution := &ContractExecution{}
	err := json.Unmarshal(msg.Data, execution)
	if err != nil {
		Log.Warningf("Failed to unmarshal contract execution during NATS tx message handling")
		nack(msg)
		return
	}

	if execution.ContractID == nil {
		Log.Errorf("Invalid tx message; missing contract_id")
		nack(msg)
		return
	}

	if execution.WalletID != nil && *execution.WalletID != uuid.Nil {
		if execution.Wallet != nil && execution.Wallet.ID != execution.Wallet.ID {
			Log.Errorf("Invalid tx message specifying a wallet_id and wallet")
			nack(msg)
			return
		}
		wallet := &Wallet{}
		wallet.setID(*execution.WalletID)
		execution.Wallet = wallet
	}

	contract := &Contract{}
	dbconf.DatabaseConnection().Where("id = ?", *execution.ContractID).Find(&contract)
	if contract == nil || contract.ID == uuid.Nil {
		Log.Errorf("Unable to execute contract; contract not found: %s", contract.ID)
		nack(msg)
		return
	}

	executionResponse, err := contract.Execute(execution)

	if err != nil {
		Log.Warningf("Failed to execute contract; %s", err.Error())
		nack(msg)
	} else {
		Log.Debugf("Executed contract; tx: %s", executionResponse)
		msg.Ack()
	}
}

func consumeTxReceiptMsg(msg *stan.Msg) {
	Log.Debugf("Consuming NATS tx receipt message: %s", msg)

	db := dbconf.DatabaseConnection()

	var tx *Transaction

	err := json.Unmarshal(msg.Data, &tx)
	if err != nil {
		desc := fmt.Sprintf("Failed to umarshal tx receipt message; %s", err.Error())
		Log.Warningf(desc)
		tx.updateStatus(db, "failed", StringOrNil(desc))
		nack(msg)
		return
	}

	tx.Reload()

	network, err := tx.GetNetwork()
	if err != nil {
		desc := fmt.Sprintf("Failed to resolve tx network; %s", err.Error())
		Log.Warningf(desc)
		tx.updateStatus(db, "failed", StringOrNil(desc))
		nack(msg)
		return
	}

	wallet, err := tx.GetWallet()
	if err != nil {
		desc := fmt.Sprintf("Failed to resolve tx wallet; %s", err.Error())
		Log.Warningf(desc)
		tx.updateStatus(db, "failed", StringOrNil(desc))
		nack(msg)
		return
	}

	err = tx.fetchReceipt(db, network, wallet)
	if err != nil {
		if msg.Redelivered { // FIXME-- implement proper dead-letter logic; only set tx to failed upon deadletter
			desc := fmt.Sprintf("Failed to fetch tx receipt; %s", err.Error())
			Log.Warningf(desc)
			tx.updateStatus(db, "failed", StringOrNil(desc))
		}

		nack(msg)
	} else {
		msg.Ack()
	}
}
