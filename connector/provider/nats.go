package provider

import (
	"encoding/json"
	"errors"
	"fmt"

	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/jinzhu/gorm"

	natsutil "github.com/kthomas/go-natsutil"
	uuid "github.com/kthomas/go.uuid"
	"github.com/provideplatform/nchain/common"
	"github.com/provideplatform/nchain/network"
	c2 "github.com/provideplatform/provide-go/api/c2"
)

// NATSProvider is a connector.ProviderAPI implementing orchestration for NATS
type NATSProvider struct {
	connectorID    uuid.UUID
	model          *gorm.DB
	config         map[string]interface{}
	networkID      *uuid.UUID
	applicationID  *uuid.UUID
	organizationID *uuid.UUID
	region         *string
	apiURL         *string
	apiPort        int
}

// InitNATSProvider initializes and returns the NATS connector API provider
func InitNATSProvider(connectorID uuid.UUID, networkID, applicationID, organizationID *uuid.UUID, model *gorm.DB, config map[string]interface{}) *NATSProvider {
	region, regionOk := config["region"].(string)
	apiURL, _ := config["api_url"].(string)
	apiPort, apiPortOk := config["api_port"].(float64)
	if connectorID == uuid.Nil || !regionOk || !apiPortOk || networkID == nil || *networkID == uuid.Nil {
		return nil
	}
	return &NATSProvider{
		connectorID:    connectorID,
		model:          model,
		config:         config,
		networkID:      networkID,
		applicationID:  applicationID,
		organizationID: organizationID,
		region:         common.StringOrNil(region),
		apiURL:         common.StringOrNil(apiURL),
		apiPort:        int(apiPort),
	}
}

func (p *NATSProvider) apiClientFactory(basePath *string) *ipfs.Shell {
	uri := ""
	if basePath != nil {
		uri = *basePath
	}
	apiURL := p.apiURLFactory(uri)
	if apiURL == nil {
		common.Log.Warningf("unable to initialize NATS api client factory")
		return nil
	}

	return ipfs.NewShell(*apiURL)
}

func (p *NATSProvider) apiURLFactory(path string) *string {
	if p.apiURL == nil {
		return nil
	}

	suffix := ""
	if path != "" {
		suffix = fmt.Sprintf("/%s", path)
	}
	return common.StringOrNil(fmt.Sprintf("%s%s", *p.apiURL, suffix))
}

func (p *NATSProvider) rawConfig() *json.RawMessage {
	cfgJSON, _ := json.Marshal(p.config)
	_cfgJSON := json.RawMessage(cfgJSON)
	return &_cfgJSON
}

// Deprovision undeploys all associated nodes and load balancers and removes them from the NATS connector
func (p *NATSProvider) Deprovision() error {
	loadBalancers := make([]*c2.LoadBalancer, 0)
	p.model.Association("LoadBalancers").Find(&loadBalancers)
	for _, balancer := range loadBalancers {
		p.model.Association("LoadBalancers").Delete(balancer)
	}

	nodes := make([]*network.Node, 0)
	p.model.Association("Nodes").Find(&nodes)
	for _, node := range nodes {
		common.Log.Debugf("Attempting to deprovision node %s on connector: %s", node.ID, p.connectorID)
		p.model.Association("Nodes").Delete(node)
		node.Delete("") // FIXME -- needs c2 API token
	}

	for _, balancer := range loadBalancers {
		msg, _ := json.Marshal(map[string]interface{}{
			"load_balancer_id": balancer.ID,
		})
		natsutil.NatsJetstreamPublish(natsLoadBalancerDeprovisioningSubject, msg)
	}

	return nil
}

// Provision configures a new load balancer and the initial NATS nodes and associates the resources with the NATS connector
func (p *NATSProvider) Provision() error {
	loadBalancer, err := c2.CreateLoadBalancer("", map[string]interface{}{
		"network_id":      *p.networkID,
		"application_id":  p.applicationID,
		"organization_id": p.organizationID,
		"type":            common.StringOrNil(ElasticsearchConnectorProvider),
		"description":     common.StringOrNil(fmt.Sprintf("NATS Connector Load Balancer")),
		"region":          p.region,
		"config":          p.config,
	})

	if err == nil {
		common.Log.Debugf("Created load balancer %s on connector: %s", loadBalancer.ID, p.connectorID)
		p.model.Association("LoadBalancers").Append(loadBalancer)

		msg, _ := json.Marshal(map[string]interface{}{
			"connector_id": p.connectorID,
		})
		natsutil.NatsJetstreamPublish(natsConnectorDenormalizeConfigSubject, msg)

		err := p.ProvisionNode()
		if err != nil {
			common.Log.Warning(err.Error())
		}
	} else {
		return fmt.Errorf("Failed to provision load balancer on connector: %s; %s", p.connectorID, err.Error())
	}

	return nil
}

// DeprovisionNode undeploys an existing node removes it from the NATS connector
func (p *NATSProvider) DeprovisionNode() error {
	node := &network.Node{}
	p.model.Association("Nodes").Find(node)

	return nil
}

// ProvisionNode deploys and load balances a new node and associates it with the NATS connector
func (p *NATSProvider) ProvisionNode() error {
	node := &network.Node{
		NetworkID:      *p.networkID,
		ApplicationID:  p.applicationID,
		OrganizationID: p.organizationID,
		Config:         p.rawConfig(),
	}

	if node.Create("") { // FIXME -- needs c2 API token
		common.Log.Debugf("Created node %s on connector: %s", node.ID, p.connectorID)
		p.model.Association("Nodes").Append(node)

		loadBalancers := make([]*c2.LoadBalancer, 0)
		p.model.Association("LoadBalancers").Find(&loadBalancers)
		for _, balancer := range loadBalancers {
			msg, _ := json.Marshal(map[string]interface{}{
				"load_balancer_id": balancer.ID.String(),
				"node_id":          node.ID.String(),
			})
			natsutil.NatsJetstreamPublish(natsLoadBalancerBalanceNodeSubject, msg)
		}

		msg, _ := json.Marshal(map[string]interface{}{
			"connector_id": p.connectorID.String(),
		})
		natsutil.NatsJetstreamPublish(natsConnectorResolveReachabilitySubject, msg)
	} else {
		return fmt.Errorf("Failed to provision node on connector: %s", p.connectorID)
	}

	return nil
}

// Reachable returns true if the NATS API provider is available
func (p *NATSProvider) Reachable() bool {
	loadBalancers := make([]*c2.LoadBalancer, 0)
	p.model.Association("LoadBalancers").Find(&loadBalancers)
	for _, loadBalancer := range loadBalancers {
		if loadBalancer.ReachableOnPort(uint(p.apiPort)) {
			return true
		}
	}
	common.Log.Debugf("Connector is unreachable: %s", p.connectorID)
	return false
}

// Create impl for NATSProvider
func (p *NATSProvider) Create(params map[string]interface{}) (*ConnectedEntity, error) {
	return nil, errors.New("create not implemented for NATS connectors")
}

// Find impl for NATSProvider
func (p *NATSProvider) Find(id string) (*ConnectedEntity, error) {
	return nil, errors.New("read not implemented for NATS connectors")
}

// Update impl for NATSProvider
func (p *NATSProvider) Update(id string, params map[string]interface{}) error {
	return errors.New("update not implemented for NATS connectors")
}

// Delete impl for NATSProvider
func (p *NATSProvider) Delete(id string) error {
	return errors.New("delete not implemented for NATS connectors")
}

// List impl for NATSProvider
func (p *NATSProvider) List(params map[string]interface{}) ([]*ConnectedEntity, error) {
	return nil, errors.New("list not implemented for NATS connectors")
}

// Query impl for NATSProvider
func (p *NATSProvider) Query(q string) (interface{}, error) {
	return nil, errors.New("query not implemented for NATS connectors")
}
