package wallet_test

import (
	"fmt"
	"math/big"
	"testing"

	bip32 "github.com/FactomProject/go-bip32"
	dbconf "github.com/kthomas/go-db-config"
	pgputil "github.com/kthomas/go-pgputil"
	uuid "github.com/kthomas/go.uuid"
	"github.com/provideapp/goldmine/common"
	"github.com/provideapp/goldmine/wallet"
	provide "github.com/provideservices/provide-go"
)

var defaultPurpose = 44

func init() {
	pgputil.RequirePGP()
}

func decrypt(w *wallet.Wallet) error {
	if w.Mnemonic != nil {
		mnemonic, err := pgputil.PGPPubDecrypt([]byte(*w.Mnemonic))
		if err != nil {
			common.Log.Warningf("Failed to decrypt mnemonic; %s", err.Error())
			return err
		}
		w.Mnemonic = common.StringOrNil(string(mnemonic))
	}

	if w.Seed != nil {
		seed, err := pgputil.PGPPubDecrypt([]byte(*w.Seed))
		if err != nil {
			common.Log.Warningf("Failed to decrypt seed; %s", err.Error())
			return err
		}
		w.Seed = common.StringOrNil(string(seed))
	}

	if w.PublicKey != nil {
		publicKey, err := pgputil.PGPPubDecrypt([]byte(*w.PublicKey))
		if err != nil {
			common.Log.Warningf("Failed to decrypt public key; %s", err.Error())
			return err
		}
		w.PublicKey = common.StringOrNil(string(publicKey))
	}

	if w.PrivateKey != nil {
		privateKey, err := pgputil.PGPPubDecrypt([]byte(*w.PrivateKey))
		if err != nil {
			common.Log.Warningf("Failed to decrypt private key; %s", err.Error())
			return err
		}
		w.PrivateKey = common.StringOrNil(string(privateKey))
	}

	return nil
}

func TestWalletCreate(t *testing.T) {
	appID, _ := uuid.NewV4()
	wallet := &wallet.Wallet{
		ApplicationID: &appID,
		Purpose:       &defaultPurpose,
	}
	if !wallet.Create() {
		t.Errorf("failed to create wallet; %s", *wallet.Errors[0].Message)
	}

	decrypt(wallet)

	masterKey, err := bip32.NewMasterKey([]byte(*wallet.Seed))
	if err != nil {
		t.Errorf("failed to init master key from seed; %s", err.Error())
	}

	masterKey2, err := bip32.NewMasterKey([]byte(*wallet.Seed))
	if err != nil {
		t.Errorf("failed to init master key from seed; %s", err.Error())
	}

	if masterKey.String() != masterKey2.String() {
		t.Errorf("failed to deterministically generate master key for seed: %s", string(*wallet.Seed))
	}

	w0, err := wallet.DeriveHardened(dbconf.DatabaseConnection(), uint32(60), uint32(0))
	if err != nil {
		t.Errorf("failed to derive hardened account; %s", err.Error())
	}
	w1, err := wallet.DeriveHardened(dbconf.DatabaseConnection(), uint32(60), uint32(0))
	if err != nil {
		t.Errorf("failed to derive hardened account; %s", err.Error())
	}

	decrypt(w0)
	decrypt(w1)

	if w0.PublicKey == nil || w1.PublicKey == nil || *w0.PublicKey != *w1.PublicKey {
		t.Errorf("failed to deterministically generate master key for seed: %s", string(*wallet.Seed))
	}

	chain := uint32(0)

	a0, err := w1.DeriveAddress(dbconf.DatabaseConnection(), uint32(0), &chain)
	if err != nil {
		t.Errorf("failed to derive address; %s", err.Error())
	}
	a1, err := w1.DeriveAddress(dbconf.DatabaseConnection(), uint32(0), &chain)
	if err != nil {
		t.Errorf("failed to derive address; %s", err.Error())
	}

	if a0.Address == "" || a1.Address == "" || a0.Address != a1.Address {
		t.Errorf("failed to deterministically generate master key for seed: %s", string(*wallet.Seed))
	}

	zerouint64 := uint64(0)

	signedTx, hash, _ := provide.EVMSignTx(
		"test-network-id",
		"http://ec2-100-27-34-141.compute-1.amazonaws.com:8050/",
		a1.Address,
		*a1.PrivateKey,
		nil,
		common.StringOrNil("0x608060405234801561001057600080fd5b50610ddc806100206000396000f3fe608060405234801561001057600080fd5b5060043610610053576000357c01000000000000000000000000000000000000000000000000000000009004806344e8fd0814610058578063ac45fcd914610074575b600080fd5b610072600480360361006d9190810190610850565b6100a4565b005b61008e600480360361008991908101906108bc565b610319565b60405161009b9190610b38565b60405180910390f35b6100ac6105ed565b6080604051908101604052803373ffffffffffffffffffffffffffffffffffffffff1681526020014281526020018481526020018381525090506000816040516020016100f99190610b98565b604051602081830303815290604052805190602001209050816001600083815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160010155604082015181600201908051906020019061019292919061062c565b5060608201518160030190805190602001906101af9291906106ac565b50905050600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190806001815401808255809150509060018203906000526020600020016000909192909190915055506000829080600181540180825580915050906001820390600052602060002090600402016000909192909190915060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506020820151816001015560408201518160020190805190602001906102b792919061062c565b5060608201518160030190805190602001906102d49291906106ac565b505050507f4cbc6aabdd0942d8df984ae683445cc9d498eff032ced24070239d9a65603bb384823360405161030b93929190610b5a565b60405180910390a150505050565b6060600083830290506000838203905060008183039050600081141561037f57600060405190808252806020026020018201604052801561037457816020015b61036161072c565b8152602001906001900390816103595790505b5093505050506105e7565b6001600080549050038311156103a057600160008054905003925081830390505b806040519080825280602002602001820160405280156103da57816020015b6103c761072c565b8152602001906001900390816103bf5790505b50935060008090505b818110156105e25760008185038154811015156103fc57fe5b9060005260206000209060040201608060405190810160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200160018201548152602001600282018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561050f5780601f106104e45761010080835404028352916020019161050f565b820191906000526020600020905b8154815290600101906020018083116104f257829003601f168201915b50505050508152602001600382018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105b15780601f10610586576101008083540402835291602001916105b1565b820191906000526020600020905b81548152906001019060200180831161059457829003601f168201915b50505050508152505085828151811015156105c857fe5b9060200190602002018190525080806001019150506103e3565b505050505b92915050565b608060405190810160405280600073ffffffffffffffffffffffffffffffffffffffff1681526020016000815260200160608152602001606081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061066d57805160ff191683800117855561069b565b8280016001018555821561069b579182015b8281111561069a57825182559160200191906001019061067f565b5b5090506106a8919061076b565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106106ed57805160ff191683800117855561071b565b8280016001018555821561071b579182015b8281111561071a5782518255916020019190600101906106ff565b5b509050610728919061076b565b5090565b608060405190810160405280600073ffffffffffffffffffffffffffffffffffffffff1681526020016000815260200160608152602001606081525090565b61078d91905b80821115610789576000816000905550600101610771565b5090565b90565b600082601f83011215156107a357600080fd5b81356107b66107b182610be7565b610bba565b915080825260208301602083018583830111156107d257600080fd5b6107dd838284610d4f565b50505092915050565b600082601f83011215156107f957600080fd5b813561080c61080782610c13565b610bba565b9150808252602083016020830185838301111561082857600080fd5b610833838284610d4f565b50505092915050565b60006108488235610d0f565b905092915050565b6000806040838503121561086357600080fd5b600083013567ffffffffffffffff81111561087d57600080fd5b610889858286016107e6565b925050602083013567ffffffffffffffff8111156108a657600080fd5b6108b285828601610790565b9150509250929050565b600080604083850312156108cf57600080fd5b60006108dd8582860161083c565b92505060206108ee8582860161083c565b9150509250929050565b60006109048383610abf565b905092915050565b61091581610d19565b82525050565b61092481610cc9565b82525050565b600061093582610c4c565b61093f8185610c85565b93508360208202850161095185610c3f565b60005b8481101561098a57838303885261096c8383516108f8565b925061097782610c78565b9150602088019750600181019050610954565b508196508694505050505092915050565b6109a481610cdb565b82525050565b60006109b582610c57565b6109bf8185610c96565b93506109cf818560208601610d5e565b6109d881610d91565b840191505092915050565b60006109ee82610c6d565b6109f88185610cb8565b9350610a08818560208601610d5e565b610a1181610d91565b840191505092915050565b6000610a2782610c62565b610a318185610ca7565b9350610a41818560208601610d5e565b610a4a81610d91565b840191505092915050565b6000608083016000830151610a6d600086018261091b565b506020830151610a806020860182610b29565b5060408301518482036040860152610a988282610a1c565b91505060608301518482036060860152610ab282826109aa565b9150508091505092915050565b6000608083016000830151610ad7600086018261091b565b506020830151610aea6020860182610b29565b5060408301518482036040860152610b028282610a1c565b91505060608301518482036060860152610b1c82826109aa565b9150508091505092915050565b610b3281610d05565b82525050565b60006020820190508181036000830152610b52818461092a565b905092915050565b60006060820190508181036000830152610b7481866109e3565b9050610b83602083018561099b565b610b90604083018461090c565b949350505050565b60006020820190508181036000830152610bb28184610a55565b905092915050565b6000604051905081810181811067ffffffffffffffff82111715610bdd57600080fd5b8060405250919050565b600067ffffff"),
		big.NewInt(0),
		&zerouint64,
		uint64(0),
	)

	if signedTx == nil {
		t.Error("failed to sign tx")
	}

	if hash == nil {
		t.Error("failed to sign tx and read resulting hash")
	}

	if signedTx != nil && hash != nil {
		fmt.Printf("signed tx... hash: %s", *hash)
	}
}
