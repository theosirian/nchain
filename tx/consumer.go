package tx

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	dbconf "github.com/kthomas/go-db-config"
	natsutil "github.com/kthomas/go-natsutil"
	uuid "github.com/kthomas/go.uuid"
	"github.com/nats-io/nats.go"
	"github.com/provideplatform/nchain/common"
	"github.com/provideplatform/nchain/contract"
	"github.com/provideplatform/nchain/wallet"
	api "github.com/provideplatform/provide-go/api"
	bookie "github.com/provideplatform/provide-go/api/bookie"
	provide "github.com/provideplatform/provide-go/api/nchain"
	util "github.com/provideplatform/provide-go/common/util"
)

const defaultNatsStream = "nchain"

// TODO: should this be calculated dynamically against average blocktime for the network and subscriptions reestablished?

const natsTxSubject = "nchain.tx"
const natsTxMaxInFlight = 1024 * 30
const natsTxMsgMaxDeliveries = 5
const txAckWait = time.Second * 60

const natsTxCreateSubject = "nchain.tx.create"
const natsTxCreateMaxInFlight = 1024 * 30
const natsTxCreateMaxDeliveries = 5
const txCreateAckWait = time.Second * 60

const natsTxFinalizeSubject = "nchain.tx.finalize"
const natsTxFinalizeMaxInFlight = 1024 * 30
const natsTxFinalizedMsgMaxDeliveries = 100
const txFinalizeAckWait = time.Second * 5

const natsTxReceiptSubject = "nchain.tx.receipt"
const natsTxReceiptMaxInFlight = 1024 * 30
const natsTxReceiptMsgMaxDeliveries = 100
const txReceiptAckWait = time.Second * 5

var waitGroup sync.WaitGroup

func init() {
	if !common.ConsumeNATSStreamingSubscriptions {
		common.Log.Debug("Tx package consumer configured to skip NATS streaming subscription setup")
		return
	}

	natsutil.EstablishSharedNatsConnection(nil)
	natsutil.NatsCreateStream(defaultNatsStream, []string{
		fmt.Sprintf("%s.>", defaultNatsStream),
	})

	createNatsTxSubscriptions(&waitGroup)
	createNatsTxCreateSubscriptions(&waitGroup)
	createNatsTxFinalizeSubscriptions(&waitGroup)
	createNatsTxReceiptSubscriptions(&waitGroup)
}

func createNatsTxSubscriptions(wg *sync.WaitGroup) {
	for i := uint64(0); i < natsutil.GetNatsConsumerConcurrency(); i++ {
		natsutil.RequireNatsJetstreamSubscription(wg,
			txAckWait,
			natsTxSubject,
			natsTxSubject,
			natsTxSubject,
			consumeTxExecutionMsg,
			txAckWait,
			natsTxMaxInFlight,
			natsTxMsgMaxDeliveries,
			nil,
		)
	}
}

func createNatsTxCreateSubscriptions(wg *sync.WaitGroup) {
	for i := uint64(0); i < natsutil.GetNatsConsumerConcurrency(); i++ {
		natsutil.RequireNatsJetstreamSubscription(wg,
			txCreateAckWait,
			natsTxCreateSubject,
			natsTxCreateSubject,
			natsTxCreateSubject,
			consumeTxCreateMsg,
			txCreateAckWait,
			natsTxCreateMaxInFlight,
			natsTxCreateMaxDeliveries,
			nil,
		)
	}
}

func createNatsTxFinalizeSubscriptions(wg *sync.WaitGroup) {
	for i := uint64(0); i < natsutil.GetNatsConsumerConcurrency(); i++ {
		natsutil.RequireNatsJetstreamSubscription(wg,
			txFinalizeAckWait,
			natsTxFinalizeSubject,
			natsTxFinalizeSubject,
			natsTxFinalizeSubject,
			consumeTxFinalizeMsg,
			txFinalizeAckWait,
			natsTxFinalizeMaxInFlight,
			natsTxFinalizedMsgMaxDeliveries,
			nil,
		)
	}
}

func createNatsTxReceiptSubscriptions(wg *sync.WaitGroup) {
	for i := uint64(0); i < natsutil.GetNatsConsumerConcurrency(); i++ {
		natsutil.RequireNatsJetstreamSubscription(wg,
			txReceiptAckWait,
			natsTxReceiptSubject,
			natsTxReceiptSubject,
			natsTxReceiptSubject,
			consumeTxReceiptMsg,
			txReceiptAckWait,
			natsTxReceiptMaxInFlight,
			natsTxReceiptMsgMaxDeliveries,
			nil,
		)
	}
}

func consumeTxCreateMsg(msg *nats.Msg) {
	common.Log.Debugf("consuming %d-byte NATS tx message on subject: %s", len(msg.Data), msg.Subject)

	var params map[string]interface{}

	err := json.Unmarshal(msg.Data, &params)
	if err != nil {
		common.Log.Warningf("failed to umarshal tx creation message; %s", err.Error())
		msg.Nak()
		return
	}

	contractID, contractIDOk := params["contract_id"]
	data, dataOk := params["data"].(string)
	accountIDStr, accountIDStrOk := params["account_id"].(string)
	walletIDStr, walletIDStrOk := params["wallet_id"].(string)
	hdDerivationPath, _ := params["hd_derivation_path"].(string)
	value, valueOk := params["value"]
	txParams, paramsOk := params["params"].(map[string]interface{})
	publishedAt, publishedAtOk := params["published_at"].(string)

	if !contractIDOk {
		common.Log.Warningf("failed to unmarshal contract_id during NATS %v message handling", msg.Subject)
		msg.Term()
		return
	}

	if !dataOk {
		common.Log.Warningf("failed to unmarshal data during NATS %v message handling", msg.Subject)
		msg.Term()
		return
	}

	if !accountIDStrOk && !walletIDStrOk {
		common.Log.Warningf("failed to unmarshal account_id or wallet_id during NATS %v message handling", msg.Subject)
		msg.Term()
		return
	}

	if !valueOk {
		common.Log.Warningf("failed to unmarshal value during NATS %v message handling", msg.Subject)
		msg.Term()
		return
	}

	if !paramsOk {
		common.Log.Warningf("failed to unmarshal params during NATS %v message handling", msg.Subject)
		msg.Term()
		return
	}

	if !publishedAtOk {
		common.Log.Warningf("failed to unmarshal published_at during NATS %v message handling", msg.Subject)
		msg.Term()
		return
	}

	contract := &contract.Contract{}
	db := dbconf.DatabaseConnection()
	db.Where("id = ?", contractID).Find(&contract)

	var accountID *uuid.UUID
	var walletID *uuid.UUID

	accountUUID, accountUUIDErr := uuid.FromString(accountIDStr)
	if accountUUIDErr == nil {
		accountID = &accountUUID
	}

	walletUUID, walletUUIDErr := uuid.FromString(walletIDStr)
	if walletUUIDErr == nil {
		walletID = &walletUUID
	}

	if accountID == nil && walletID == nil {
		common.Log.Warningf("failed to unmarshal account_id or wallet_id during NATS %v message handling", msg.Subject)
		msg.Term()
		return
	}

	publishedAtTime, err := time.Parse(time.RFC3339, publishedAt)
	if err != nil {
		common.Log.Warningf("failed to parse published_at as RFC3339 timestamp during NATS %v message handling; %s", msg.Subject, err.Error())
		msg.Term()
		return
	}

	tx := &Transaction{
		ApplicationID:  contract.ApplicationID,
		OrganizationID: contract.OrganizationID,
		Data:           common.StringOrNil(data),
		NetworkID:      contract.NetworkID,
		AccountID:      accountID,
		WalletID:       walletID,
		Path:           common.StringOrNil(hdDerivationPath),
		To:             nil,
		Value:          &TxValue{value: big.NewInt(int64(value.(float64)))},
		PublishedAt:    &publishedAtTime,
	}
	tx.setParams(txParams)

	if tx.Create(db) {
		contract.TransactionID = &tx.ID
		db.Save(&contract)
		common.Log.Debugf("transaction execution successful: %s", *tx.Hash)
		msg.Ack()
	} else {
		errmsg := fmt.Sprintf("failed to execute transaction; tx failed with %d error(s)", len(tx.Errors))
		for _, err := range tx.Errors {
			errmsg = fmt.Sprintf("%s\n\t%s", errmsg, *err.Message)
		}

		faucetSubsidyEligible := strings.Contains(strings.ToLower(errmsg), "insufficient funds") && tx.shouldSubsidize()
		if faucetSubsidyEligible {
			common.Log.Debugf("transaction execution failed due to insufficient funds but faucet subsidy exists for network: %s; requesting subsidized tx funding", tx.NetworkID)
			faucetBeneficiary, _ := tx.signerFactory(db)
			faucetBeneficiaryAddress := faucetBeneficiary.Address()

			params := tx.ParseParams()
			gas, gasOk := params["gas"].(float64)
			if !gasOk {
				gas = float64(210000 * 2) // FIXME-- parameterize
			}
			networkSubsidyFaucetDripValue := int64(100000000000000000) // FIXME-- configurable

			err = subsidize(
				db,
				tx.NetworkID,
				faucetBeneficiaryAddress,
				networkSubsidyFaucetDripValue,
				int64(gas),
			)
			if err == nil {
				db.Delete(&tx) // Drop tx that had insufficient funds so its hash can be rebroadcast...
				common.Log.Debugf("faucet subsidy transaction broadcast; beneficiary: %s", faucetBeneficiaryAddress)
				contract.TransactionID = &tx.ID
				db.Save(&contract)
				msg.Ack()
				return
			}
		} else {
			common.Log.Warning(errmsg)
			msg.Nak()
		}
	}
}

// subsidize the given beneficiary with a drip equal to the given val
func subsidize(db *gorm.DB, networkID uuid.UUID, beneficiary string, val, gas int64) error {
	payment, err := bookie.CreatePayment(util.DefaultVaultAccessJWT, map[string]interface{}{
		"to":   beneficiary,
		"data": common.StringOrNil("0x"),
	})
	if err != nil {
		return err
	}

	common.Log.Debugf("subsidized transaction using api.providepayments.com; beneficiary: %s; tx hash: %s", beneficiary, payment.Params["result"].(string))
	return nil
}

func consumeTxExecutionMsg(msg *nats.Msg) {
	common.Log.Debugf("consuming %d-byte NATS tx message on subject: %s", len(msg.Data), msg.Subject)

	execution := &contract.Execution{}
	err := json.Unmarshal(msg.Data, execution)
	if err != nil {
		common.Log.Warningf("failed to unmarshal contract execution during NATS tx message handling")
		msg.Nak()
		return
	}

	if execution.ContractID == nil {
		common.Log.Errorf("invalid tx message; missing contract_id")
		msg.Term()
		return
	}

	if execution.AccountID != nil && *execution.AccountID != uuid.Nil {
		var executionAccountID *uuid.UUID
		if executionAccount, executionAccountOk := execution.Account.(map[string]interface{}); executionAccountOk {
			if executionAccountIDStr, executionAccountIDStrOk := executionAccount["id"].(string); executionAccountIDStrOk {
				execAccountUUID, err := uuid.FromString(executionAccountIDStr)
				if err == nil {
					executionAccountID = &execAccountUUID
				}
			}
		}
		if execution.Account != nil && execution.AccountID != nil && *executionAccountID != *execution.AccountID {
			common.Log.Errorf("invalid tx message specifying a account_id and account")
			msg.Term()
			return
		}
		account := &wallet.Account{}
		account.SetID(*execution.AccountID)
		execution.Account = account
	}

	if execution.WalletID != nil && *execution.WalletID != uuid.Nil {
		var executionWalletID *uuid.UUID
		if executionWallet, executionWalletOk := execution.Wallet.(map[string]interface{}); executionWalletOk {
			if executionWalletIDStr, executionWalletIDStrOk := executionWallet["id"].(string); executionWalletIDStrOk {
				execWalletUUID, err := uuid.FromString(executionWalletIDStr)
				if err == nil {
					executionWalletID = &execWalletUUID
				}
			}
		}
		if execution.Wallet != nil && execution.WalletID != nil && *executionWalletID != *execution.WalletID {
			common.Log.Errorf("invalid tx message specifying a wallet_id and wallet")
			msg.Term()
			return
		}
		wallet := &wallet.Wallet{}
		wallet.SetID(*execution.WalletID)
		execution.Wallet = wallet
	}

	db := dbconf.DatabaseConnection()

	cntract := &contract.Contract{}
	db.Where("id = ?", *execution.ContractID).Find(&cntract)
	if cntract == nil || cntract.ID == uuid.Nil {
		db.Where("address = ?", *execution.ContractID).Find(&cntract)
	}
	if cntract == nil || cntract.ID == uuid.Nil {
		common.Log.Errorf("unable to execute contract; contract not found: %s", cntract.ID)
		msg.Term()
		return
	}

	executionResponse, err := executeTransaction(cntract, execution)
	if err != nil {
		common.Log.Debugf("contract execution failed; %s", err.Error())

		// CHECKME - this functionality is now in bookie, and shouldn't be replicated here
		// if execution.AccountAddress != nil {
		// 	networkSubsidyFaucetDripValue := int64(100000000000000000) // FIXME-- configurable
		// 	_subsidize := strings.Contains(strings.ToLower(err.Error()), "insufficient funds") && tx.shouldSubsidize()

		// 	if _subsidize {
		// 		common.Log.Debugf("contract execution failed due to insufficient funds but tx subsidize flag is set; requesting subsidized tx funding for target network: %s", cntract.NetworkID)
		// 		faucetBeneficiaryAddress := *execution.AccountAddress
		// 		err = subsidize(
		// 			db,
		// 			cntract.NetworkID,
		// 			faucetBeneficiaryAddress,
		// 			networkSubsidyFaucetDripValue,
		// 			int64(210000*2),
		// 		)
		// 		if err == nil {
		// 			db.Where("ref = ?", execution.Ref).Find(&tx)
		// 			if tx != nil && tx.ID != uuid.Nil {
		// 				db.Delete(&tx) // Drop tx that had insufficient funds so its hash can be rebroadcast...
		// 			}
		// 			common.Log.Debugf("faucet subsidy transaction broadcast; beneficiary: %s", faucetBeneficiaryAddress)
		// 		}
		// 	}

		// } else {
		// 	common.Log.Warningf("Failed to execute contract; %s", err.Error())
		// }

		msg.Nak()
	} else {
		logmsg := fmt.Sprintf("executed contract: %s", *cntract.Address)
		if executionResponse != nil && executionResponse.Response != nil {
			logmsg = fmt.Sprintf("%s; response: %s", logmsg, executionResponse.Response)
		}
		common.Log.Debug(logmsg)

		msg.Ack()
	}
}

func consumeTxFinalizeMsg(msg *nats.Msg) {
	common.Log.Tracef("consuming NATS tx finalize message: %s", msg)

	var params map[string]interface{}

	nack := func(msg *nats.Msg, errmsg string, dropPacket bool) {
		if dropPacket {
			common.Log.Tracef("dropping tx packet on the floor; %s", errmsg)
			msg.Term()
			return
		}
		msg.Nak()
	}

	err := json.Unmarshal(msg.Data, &params)
	if err != nil {
		nack(msg, fmt.Sprintf("failed to umarshal tx finalize message; %s", err.Error()), true)
		return
	}

	block, blockOk := params["block"].(float64)
	blockTimestampStr, blockTimestampStrOk := params["block_timestamp"].(string)
	finalizedAtStr, finalizedAtStrOk := params["finalized_at"].(string)
	hash, hashOk := params["hash"].(string)

	if !blockOk {
		nack(msg, "failed to finalize tx; no block provided", true)
		return
	}

	if !blockTimestampStrOk {
		nack(msg, "failed to finalize tx; no block timestamp provided", true)
		return
	}

	if !finalizedAtStrOk {
		nack(msg, "failed to finalize tx; no finalized at timestamp provided", true)
		return
	}

	if !hashOk {
		nack(msg, "failed to finalize tx; no hash provided", true)
		return
	}

	blockTimestamp, err := time.Parse(time.RFC3339, blockTimestampStr)
	if err != nil {
		nack(msg, fmt.Sprintf("failed to unmarshal block_timestamp during NATS %v message handling; %s", msg.Subject, err.Error()), true)
		return
	}

	finalizedAt, err := time.Parse(time.RFC3339, finalizedAtStr)
	if err != nil {
		nack(msg, fmt.Sprintf("failed to unmarshal finalized_at during NATS %v message handling; %s", msg.Subject, err.Error()), true)
		return
	}

	tx := &Transaction{}
	db := dbconf.DatabaseConnection()

	common.Log.Tracef("checking local db for tx status; tx hash: %s", hash)

	db.Where("hash = ? AND status IN (?, ?)", hash, "pending", "failed").Find(&tx)
	if tx == nil || tx.ID == uuid.Nil {
		// TODO: this is integration point to upsert Wallet & Transaction... need to think thru performance implications & implementation details
		nack(msg, fmt.Sprintf("failed to mark block and finalized_at timestamp on tx: %s; tx not found for given hash", hash), true)
		return
	}

	blockNumber := uint64(block)

	tx.Block = &blockNumber
	tx.BlockTimestamp = &blockTimestamp
	tx.FinalizedAt = &finalizedAt
	if tx.BroadcastAt != nil {
		if tx.PublishedAt != nil {
			queueLatency := uint64(tx.BroadcastAt.Sub(*tx.PublishedAt)) / uint64(time.Millisecond)
			tx.QueueLatency = &queueLatency

			e2eLatency := uint64(tx.FinalizedAt.Sub(*tx.PublishedAt)) / uint64(time.Millisecond)
			tx.E2ELatency = &e2eLatency
		}

		networkLatency := uint64(tx.FinalizedAt.Sub(*tx.BroadcastAt)) / uint64(time.Millisecond)
		tx.NetworkLatency = &networkLatency
	}

	tx.updateStatus(db, "success", nil)
	result := db.Save(&tx)
	errors := result.GetErrors()
	if len(errors) > 0 {
		for _, err := range errors {
			tx.Errors = append(tx.Errors, &api.Error{
				Message: common.StringOrNil(err.Error()),
			})
		}
	}
	if len(tx.Errors) > 0 {
		nack(msg, fmt.Sprintf("Failed to set block and finalized_at timestamp on tx: %s; error: %s", hash, *tx.Errors[0].Message), false)
		return
	}

	msg.Ack()
}

func consumeTxReceiptMsg(msg *nats.Msg) {
	defer func() {
		if r := recover(); r != nil {
			common.Log.Warningf("recovered from failed tx receipt message; %s", r)
			msg.Term()
		}
	}()

	common.Log.Debugf("consuming NATS tx receipt message: %s", msg)

	var params map[string]interface{}

	err := json.Unmarshal(msg.Data, &params)
	if err != nil {
		common.Log.Warningf("failed to umarshal load balancer provisioning message; %s", err.Error())
		msg.Nak()
		return
	}

	transactionID, transactionIDOk := params["transaction_id"].(string)
	if !transactionIDOk {
		common.Log.Warningf("failed to consume NATS tx receipt message; no transaction id provided")
		msg.Term()
		return
	}

	db := dbconf.DatabaseConnection()

	tx := &Transaction{}
	db.Where("id = ?", transactionID).Find(&tx)
	if tx == nil || tx.ID == uuid.Nil {
		common.Log.Tracef("failed to fetch tx receipt; no tx resolved for id: %s", transactionID)
		msg.Term()
		return
	}

	signer, err := tx.signerFactory(db)
	if err != nil {
		desc := "failed to resolve tx signing account or HD wallet"
		common.Log.Warningf(desc)
		tx.updateStatus(db, "failed", common.StringOrNil(desc))
		msg.Nak()
		return
	}

	err = tx.fetchReceipt(db, signer.Network, signer.Address())
	if err != nil {
		common.Log.Debugf(fmt.Sprintf("failed to fetch tx receipt; %s", err.Error()))
		// msg.Nak()
		return
	} else {
		common.Log.Debugf("fetched tx receipt for hash: %s", *tx.Hash)

		blockNumber := tx.Response.Receipt.(*provide.TxReceipt).BlockNumber
		// if we have a block number in the receipt, and the tx has no block
		// populate the block and finalized timestamp
		if blockNumber != nil && tx.Block == nil {
			receiptBlock := blockNumber.Uint64()
			tx.Block = &receiptBlock
			receiptFinalized := time.Now()
			tx.FinalizedAt = &receiptFinalized
			common.Log.Debugf("tx %s finalized in block %v at %s", *tx.Hash, blockNumber, receiptFinalized.Format("Mon, 02 Jan 2006 15:04:05 MST"))
		}

		tx.updateStatus(db, "success", nil)
		msg.Ack()
	}
}
