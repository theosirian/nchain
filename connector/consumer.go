package connector

import (
	"encoding/json"
	"sync"
	"time"

	dbconf "github.com/kthomas/go-db-config"
	natsutil "github.com/kthomas/go-natsutil"
	uuid "github.com/kthomas/go.uuid"
	stan "github.com/nats-io/stan.go"
	"github.com/provideapp/goldmine/common"
)

const natsConnectorDeprovisioningSubject = "goldmine.connector.provision"
const natsConnectorDeprovisioningMaxInFlight = 64
const natsConnectorDeprovisioningTimeout = int64(time.Minute * 10)
const natsConnectorDeprovisioningInvocationTimeout = time.Second * 15

const natsConnectorProvisioningSubject = "goldmine.connector.deprovision"
const natsConnectorProvisioningMaxInFlight = 64
const natsConnectorProvisioningTimeout = int64(time.Minute * 10)
const natsConnectorProvisioningInvocationTimeout = time.Second * 15

const natsConnectorDenormalizeConfigSubject = "goldmine.connector.config.denormalize"
const natsConnectorDenormalizeConfigMaxInFlight = 64
const natsConnectorDenormalizeConfigTimeout = int64(time.Minute * 1)
const natsConnectorDenormalizeConfigInvocationTimeout = time.Second * 5

var waitGroup sync.WaitGroup

func init() {
	if !common.ConsumeNATSStreamingSubscriptions {
		common.Log.Debug("Connector package consumer configured to skip NATS streaming subscription setup")
		return
	}

	createNatsConnectorProvisioningSubscriptions(&waitGroup)
	createNatsConnectorDeprovisioningSubscriptions(&waitGroup)
	createNatsConnectorDenormalizeConfigSubscriptions(&waitGroup)
}

func createNatsConnectorProvisioningSubscriptions(wg *sync.WaitGroup) {
	for i := uint64(0); i < natsutil.GetNatsConsumerConcurrency(); i++ {
		natsutil.RequireNatsStreamingSubscription(wg,
			natsConnectorProvisioningInvocationTimeout,
			natsConnectorProvisioningSubject,
			natsConnectorProvisioningSubject,
			consumeConnectorProvisioningMsg,
			natsConnectorProvisioningInvocationTimeout,
			natsConnectorProvisioningMaxInFlight,
		)
	}
}

func createNatsConnectorDeprovisioningSubscriptions(wg *sync.WaitGroup) {
	for i := uint64(0); i < natsutil.GetNatsConsumerConcurrency(); i++ {
		natsutil.RequireNatsStreamingSubscription(wg,
			natsConnectorDeprovisioningInvocationTimeout,
			natsConnectorDeprovisioningSubject,
			natsConnectorDeprovisioningSubject,
			consumeConnectorDeprovisioningMsg,
			natsConnectorDeprovisioningInvocationTimeout,
			natsConnectorDeprovisioningMaxInFlight,
		)
	}
}

func createNatsConnectorDenormalizeConfigSubscriptions(wg *sync.WaitGroup) {
	for i := uint64(0); i < natsutil.GetNatsConsumerConcurrency(); i++ {
		natsutil.RequireNatsStreamingSubscription(wg,
			natsConnectorDenormalizeConfigInvocationTimeout,
			natsConnectorDenormalizeConfigSubject,
			natsConnectorDenormalizeConfigSubject,
			consumeConnectorDenormalizeConfigMsg,
			natsConnectorDenormalizeConfigInvocationTimeout,
			natsConnectorDenormalizeConfigMaxInFlight,
		)
	}
}

func consumeConnectorProvisioningMsg(msg *stan.Msg) {
	common.Log.Debugf("Consuming NATS connector provisioning message: %s", msg)
	var params map[string]interface{}

	err := json.Unmarshal(msg.Data, &params)
	if err != nil {
		common.Log.Warningf("Failed to umarshal connector provisioning message; %s", err.Error())
		natsutil.Nack(msg)
		return
	}

	connectorID, connectorIDOk := params["connector_id"].(string)
	if !connectorIDOk {
		common.Log.Warningf("Failed to provision connector; no connector id provided")
		natsutil.Nack(msg)
		return
	}

	db := dbconf.DatabaseConnection()

	connector := &Connector{}
	db.Where("id = ?", connectorID).Find(&connector)
	if connector == nil || connector.ID == uuid.Nil {
		common.Log.Warningf("Failed to provision connector; no connector resolved for id: %s", connectorID)
		natsutil.Nack(msg)
		return
	}

	err = connector.provision()
	if err != nil {
		common.Log.Warningf("Failed to provision connector; %s", err.Error())
		natsutil.AttemptNack(msg, natsConnectorProvisioningTimeout)
	} else {
		common.Log.Debugf("Connector provisioning succeeded; ACKing NATS message for connector: %s", connector.ID)
		msg.Ack()
	}
}

func consumeConnectorDeprovisioningMsg(msg *stan.Msg) {
	common.Log.Debugf("Consuming NATS connector deprovisioning message: %s", msg)
	var params map[string]interface{}

	err := json.Unmarshal(msg.Data, &params)
	if err != nil {
		common.Log.Warningf("Failed to umarshal connector deprovisioning message; %s", err.Error())
		natsutil.Nack(msg)
		return
	}

	connectorID, connectorIDOk := params["connector_id"].(string)
	if !connectorIDOk {
		common.Log.Warningf("Failed to deprovision connector; no connector id provided")
		natsutil.Nack(msg)
		return
	}

	db := dbconf.DatabaseConnection()

	connector := &Connector{}
	db.Where("id = ?", connectorID).Find(&connector)
	if connector == nil || connector.ID == uuid.Nil {
		common.Log.Warningf("Failed to deprovision connector; no connector resolved for id: %s", connectorID)
		natsutil.Nack(msg)
		return
	}

	err = connector.deprovision()
	if err != nil {
		common.Log.Warningf("Failed to deprovision connector; %s", err.Error())
		natsutil.AttemptNack(msg, natsConnectorDeprovisioningTimeout)
	} else {
		common.Log.Debugf("Connector deprovisioning succeeded; ACKing NATS message for connector: %s", connector.ID)
		msg.Ack()
	}
}

func consumeConnectorDenormalizeConfigMsg(msg *stan.Msg) {
	common.Log.Debugf("Consuming NATS connector denormalize config message: %s", msg)
	var params map[string]interface{}

	err := json.Unmarshal(msg.Data, &params)
	if err != nil {
		common.Log.Warningf("Failed to umarshal connector denormalize config message; %s", err.Error())
		natsutil.Nack(msg)
		return
	}

	connectorID, connectorIDOk := params["connector_id"].(string)
	if !connectorIDOk {
		common.Log.Warningf("Failed to denormalize connector config; no connector id provided")
		natsutil.Nack(msg)
		return
	}

	db := dbconf.DatabaseConnection()

	connector := &Connector{}
	db.Where("id = ?", connectorID).Find(&connector)
	if connector == nil || connector.ID == uuid.Nil {
		common.Log.Warningf("Failed to denormalize connector config; no connector resolved for id: %s", connectorID)
		natsutil.Nack(msg)
		return
	}

	err = connector.denormalizeConfig()
	if err != nil {
		common.Log.Warningf("Failed to denormalize connector config; %s", err.Error())
		natsutil.AttemptNack(msg, natsConnectorDenormalizeConfigTimeout)
	} else {
		common.Log.Debugf("Connector config denormalized; ACKing NATS message for connector: %s", connector.ID)
		msg.Ack()
	}
}
