package provider

import (
	"encoding/json"
	"errors"
	"fmt"

	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/jinzhu/gorm"

	natsutil "github.com/kthomas/go-natsutil"
	uuid "github.com/kthomas/go.uuid"
	"github.com/provideapp/goldmine/common"
	"github.com/provideapp/goldmine/network"
)

// TableauProvider is a connector.ProviderAPI implementing orchestration for Tableau
type TableauProvider struct {
	connectorID    uuid.UUID
	model          *gorm.DB
	config         map[string]interface{}
	networkID      *uuid.UUID
	applicationID  *uuid.UUID
	organizationID *uuid.UUID
	region         *string
	apiURL         *string
	apiPort        int
	gatewayPort    int
}

// InitTableauProvider initializes and returns the Tableau connector API provider
func InitTableauProvider(connectorID uuid.UUID, networkID, applicationID, organizationID *uuid.UUID, model *gorm.DB, config map[string]interface{}) *TableauProvider {
	region, regionOk := config["region"].(string)
	apiURL, _ := config["api_url"].(string)
	apiPort, apiPortOk := config["api_port"].(float64)
	gatewayPort, gatewayPortOk := config["gateway_port"].(float64)
	if connectorID == uuid.Nil || !regionOk || !apiPortOk || !gatewayPortOk || networkID == nil || *networkID == uuid.Nil {
		return nil
	}
	return &TableauProvider{
		connectorID:    connectorID,
		model:          model,
		config:         config,
		networkID:      networkID,
		applicationID:  applicationID,
		organizationID: organizationID,
		region:         common.StringOrNil(region),
		apiURL:         common.StringOrNil(apiURL),
		apiPort:        int(apiPort),
		gatewayPort:    int(gatewayPort),
	}
}

func (p *TableauProvider) apiClientFactory(basePath *string) *ipfs.Shell {
	uri := ""
	if basePath != nil {
		uri = *basePath
	}
	apiURL := p.apiURLFactory(uri)
	if apiURL == nil {
		common.Log.Warningf("unable to initialize Tableau api client factory")
		return nil
	}

	return ipfs.NewShell(*apiURL)
}

func (p *TableauProvider) apiURLFactory(path string) *string {
	if p.apiURL == nil {
		return nil
	}

	suffix := ""
	if path != "" {
		suffix = fmt.Sprintf("/%s", path)
	}
	return common.StringOrNil(fmt.Sprintf("%s%s", *p.apiURL, suffix))
}

func (p *TableauProvider) rawConfig() *json.RawMessage {
	cfgJSON, _ := json.Marshal(p.config)
	_cfgJSON := json.RawMessage(cfgJSON)
	return &_cfgJSON
}

// Deprovision undeploys all associated nodes and load balancers and removes them from the Tableau connector
func (p *TableauProvider) Deprovision() error {
	loadBalancers := make([]*network.LoadBalancer, 0)
	p.model.Association("LoadBalancers").Find(&loadBalancers)
	for _, balancer := range loadBalancers {
		p.model.Association("LoadBalancers").Delete(balancer)
	}

	nodes := make([]*network.Node, 0)
	p.model.Association("Nodes").Find(&nodes)
	for _, node := range nodes {
		common.Log.Debugf("Attempting to deprovision node %s on connector: %s", node.ID, p.connectorID)
		p.model.Association("Nodes").Delete(node)
		node.Delete()
	}

	for _, balancer := range loadBalancers {
		msg, _ := json.Marshal(map[string]interface{}{
			"load_balancer_id": balancer.ID,
		})
		natsutil.NatsStreamingPublish(natsLoadBalancerDeprovisioningSubject, msg)
	}

	return nil
}

// Provision configures a new load balancer and the initial Tableau nodes and associates the resources with the Tableau connector
func (p *TableauProvider) Provision() error {
	loadBalancer := &network.LoadBalancer{
		NetworkID:      *p.networkID,
		ApplicationID:  p.applicationID,
		OrganizationID: p.organizationID,
		Type:           common.StringOrNil(TableauConnectorProvider),
		Description:    common.StringOrNil(fmt.Sprintf("Tableau Connector Load Balancer")),
		Region:         p.region,
		Config:         p.rawConfig(),
	}

	if loadBalancer.Create() {
		common.Log.Debugf("Created load balancer %s on connector: %s", loadBalancer.ID, p.connectorID)
		p.model.Association("LoadBalancers").Append(loadBalancer)

		msg, _ := json.Marshal(map[string]interface{}{
			"connector_id": p.connectorID,
		})
		natsutil.NatsStreamingPublish(natsConnectorDenormalizeConfigSubject, msg)

		err := p.ProvisionNode()
		if err != nil {
			common.Log.Warning(err.Error())
		}
	} else {
		return fmt.Errorf("Failed to provision load balancer on connector: %s; %s", p.connectorID, *loadBalancer.Errors[0].Message)
	}

	return nil
}

// DeprovisionNode undeploys an existing node removes it from the Tableau connector
func (p *TableauProvider) DeprovisionNode() error {
	node := &network.Node{}
	p.model.Association("Nodes").Find(node)

	return nil
}

// ProvisionNode deploys and load balances a new node and associates it with the Tableau connector
func (p *TableauProvider) ProvisionNode() error {
	node := &network.Node{
		NetworkID:      *p.networkID,
		ApplicationID:  p.applicationID,
		OrganizationID: p.organizationID,
		Config:         p.rawConfig(),
	}

	if node.Create() {
		common.Log.Debugf("Created node %s on connector: %s", node.ID, p.connectorID)
		p.model.Association("Nodes").Append(node)

		loadBalancers := make([]*network.LoadBalancer, 0)
		p.model.Association("LoadBalancers").Find(&loadBalancers)
		for _, balancer := range loadBalancers {
			msg, _ := json.Marshal(map[string]interface{}{
				"load_balancer_id": balancer.ID.String(),
				"network_node_id":  node.ID.String(),
			})
			natsutil.NatsStreamingPublish(natsLoadBalancerBalanceNodeSubject, msg)
		}

		msg, _ := json.Marshal(map[string]interface{}{
			"connector_id": p.connectorID.String(),
		})
		natsutil.NatsStreamingPublish(natsConnectorResolveReachabilitySubject, msg)
	} else {
		return fmt.Errorf("Failed to provision node on connector: %s", p.connectorID)
	}

	return nil
}

// Reachable returns true if the Tableau API provider is available
func (p *TableauProvider) Reachable() bool {
	loadBalancers := make([]*network.LoadBalancer, 0)
	p.model.Association("LoadBalancers").Find(&loadBalancers)
	for _, loadBalancer := range loadBalancers {
		if loadBalancer.ReachableOnPort(uint(p.apiPort)) {
			return true
		}
	}
	common.Log.Debugf("Connector is unreachable: %s", p.connectorID)
	return false
}

// Create impl for TableauProvider
func (p *TableauProvider) Create(params map[string]interface{}) (*ConnectedEntity, error) {
	return nil, errors.New("create not implemented for Tableau connectors")
}

// Read impl for TableauProvider
func (p *TableauProvider) Read(id string) (*ConnectedEntity, error) {
	return nil, errors.New("read not implemented for Tableau connectors")
}

// Update impl for TableauProvider
func (p *TableauProvider) Update(id string, params map[string]interface{}) error {
	return errors.New("update not implemented for Tableau connectors")
}

// Delete impl for TableauProvider
func (p *TableauProvider) Delete(id string) error {
	return errors.New("delete not implemented for Tableau connectors")
}

// List impl for TableauProvider
func (p *TableauProvider) List(params map[string]interface{}) ([]*ConnectedEntity, error) {
	return nil, errors.New("list not implemented for Tableau connectors")
}

// Query impl for TableauProvider
func (p *TableauProvider) Query(q string) (interface{}, error) {
	return nil, errors.New("query not implemented for Tableau connectors")
}
