package p2p

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/kthomas/go.uuid"
	"github.com/provideapp/goldmine/common"
)

// PlatformBcoin bcoin platform
const PlatformBcoin = "bcoin"

// PlatformEVM evm platform
const PlatformEVM = "evm"

// PlatformHandshake handshake platform
const PlatformHandshake = "handshake"

// PlatformHyperledger hyperledger platform
const PlatformHyperledger = "hyperledger"

// PlatformQuorum quorum platform
const PlatformQuorum = "quorum"

// ProviderBcoin bcoin p2p provider
const ProviderBcoin = "bcoin"

// ProviderGeth geth p2p provider
const ProviderGeth = "geth"

// ProviderParity parity p2p provider
const ProviderParity = "parity"

// ProviderQuorum quorum p2p provider
const ProviderQuorum = "quorum"

// API defines an interface for p2p network implementations
type API interface {
	AcceptNonReservedPeers() error
	DropNonReservedPeers() error
	AddPeer(string) error
	RemovePeer(string) error
	ParsePeerURL(string) (*string, error)
	FormatBootnodes([]string) string
	RequireBootnodes(db *gorm.DB, userID *uuid.UUID, networkID *uuid.UUID, n common.Configurable) error
	ResolvePeerURL() (*string, error)
	Upgrade() error

	DefaultEntrypoint() []string
	EnrichStartCommand() []string
}
