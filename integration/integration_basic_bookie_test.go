// +build integration nchain bookie basic

package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	uuid "github.com/kthomas/go.uuid"
	"github.com/provideplatform/ident/common"
	nchain "github.com/provideplatform/provide-go/api/nchain"
)

func TestContractWalletNoDerivationPath(t *testing.T) {

	t.Parallel()

	testId, err := uuid.NewV4()
	if err != nil {
		t.Logf("error creating new UUID")
	}

	userToken, err := UserAndTokenFactory(testId)
	if err != nil {
		t.Errorf("user authentication failed. Error: %s", err.Error())
	}

	testcaseApp := Application{
		"app" + testId.String(),
		"appdesc " + testId.String(),
	}

	app, err := appFactory(*userToken, testcaseApp.name, testcaseApp.description)
	if err != nil {
		t.Errorf("error setting up application. Error: %s", err.Error())
		return
	}

	appToken, err := appTokenFactory(*userToken, app.ID)
	if err != nil {
		t.Errorf("error getting app token. Error: %s", err.Error())
		return
	}

	wallet, err := nchain.CreateWallet(*appToken.Token, map[string]interface{}{
		"mnemonic": "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err != nil {
		t.Errorf("error creating wallet: %s", err.Error())
		return
	}

	// load the ekho compiled artifact
	ekhoArtifact, err := ioutil.ReadFile("artifacts/ekho.json")
	if err != nil {
		t.Errorf("error loading ekho artifact. Error: %s", err.Error())
	}

	ekhoCompiledArtifact := nchain.CompiledArtifact{}
	err = json.Unmarshal(ekhoArtifact, &ekhoCompiledArtifact)
	if err != nil {
		t.Errorf("error converting ekho compiled artifact. Error: %s", err.Error())
	}

	tt := []struct {
		network  string
		name     string
		walletID string
		artifact nchain.CompiledArtifact
	}{
		{rinkebyNetworkID, "ekho", wallet.ID.String(), ekhoCompiledArtifact},
		//{ropstenNetworkID, "ekho", wallet.ID.String(), ekhoCompiledArtifact},
		//{kovanNetworkID, "ekho", wallet.ID.String(), ekhoCompiledArtifact},
		//{goerliNetworkID, "ekho", wallet.ID.String(), ekhoCompiledArtifact}, // note: goerli network has provide-go bug, likely present in bookie as well
	}

	for _, tc := range tt {

		t.Logf("creating contract using wallet id: %s", tc.walletID)
		contract, err := nchain.CreateContract(*appToken.Token, map[string]interface{}{
			"network_id":     tc.network,
			"application_id": app.ID.String(),
			"wallet_id":      tc.walletID,
			"name":           tc.name,
			"address":        "0x",
			"params": map[string]interface{}{
				"wallet_id":         tc.walletID,
				"compiled_artifact": tc.artifact,
			},
		})
		if err != nil {
			t.Errorf("error creating %s contract. Error: %s", tc.name, err.Error())
			return
		}

		// wait for the contract to be deployed
		started := time.Now().Unix()
		for {
			if time.Now().Unix()-started >= contractTimeout {
				t.Errorf("timed out awaiting contract address for %s contract for %s network", tc.name, tc.network)
				return
			}

			cntrct, err := nchain.GetContractDetails(*appToken.Token, contract.ID.String(), map[string]interface{}{})
			if err != nil {
				t.Errorf("error fetching %s contract details; %s", tc.name, err.Error())
				return
			}

			if cntrct.Address != nil && *cntrct.Address != "0x" {
				t.Logf("%s contract address resolved; contract id: %s; address: %s", tc.name, cntrct.ID.String(), *cntrct.Address)
				break
			}

			t.Logf("resolving contract for %s network ...", tc.network)
			time.Sleep(contractSleepTime * time.Second)
		}

		// create a message for ekho
		msg := common.RandomString(118)
		t.Logf("executing contract using wallet id: %s, nil derivation path", tc.walletID)

		params := map[string]interface{}{}

		parameter := fmt.Sprintf(`{"method":"broadcast", "params": ["%s"], "value":0, "wallet_id":"%s"}`, msg, tc.walletID)
		json.Unmarshal([]byte(parameter), &params)

		// execute the contract method
		execResponse, err := nchain.ExecuteContract(*appToken.Token, contract.ID.String(), params)
		if err != nil {
			t.Errorf("error executing contract. Error: %s", err.Error())
			return
		}

		if err != nil {
			t.Errorf("error executing contract: %s", err.Error())
			return
		}

		// wait for the transaction to be mined (get a tx hash)
		started = time.Now().Unix()
		for {
			if time.Now().Unix()-started >= transactionTimeout {
				t.Error("timed out awaiting transaction hash")
				return
			}

			tx, err := nchain.GetTransactionDetails(*appToken.Token, *execResponse.Reference, map[string]interface{}{})
			//this is populated by nchain consumer, so it can take a moment to appear, so we won't quit right away on a 404
			if err != nil {
				t.Logf("tx not yet mined...")
			}

			if err == nil {
				if tx.Hash != nil && *tx.Hash != "0x" {
					t.Logf("tx resolved; tx id: %s; hash: %s", tx.ID.String(), *tx.Hash)
					break
				}
				t.Logf("resolving transaction...")
			}
			time.Sleep(transactionSleepTime * time.Second)
		}

		t.Logf("contract execution successful")
	}
}

func TestContractWallet(t *testing.T) {

	t.Parallel()

	testId, err := uuid.NewV4()
	if err != nil {
		t.Logf("error creating new UUID")
	}

	userToken, err := UserAndTokenFactory(testId)
	if err != nil {
		t.Errorf("user authentication failed. Error: %s", err.Error())
	}

	testcaseApp := Application{
		"app" + testId.String(),
		"appdesc " + testId.String(),
	}

	app, err := appFactory(*userToken, testcaseApp.name, testcaseApp.description)
	if err != nil {
		t.Errorf("error setting up application. Error: %s", err.Error())
		return
	}

	appToken, err := appTokenFactory(*userToken, app.ID)
	if err != nil {
		t.Errorf("error getting app token. Error: %s", err.Error())
		return
	}

	wallet, err := nchain.CreateWallet(*appToken.Token, map[string]interface{}{
		"mnemonic": "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err != nil {
		t.Errorf("error creating wallet: %s", err.Error())
		return
	}

	// this path produces a 0-balance eth address, forcing bookie to subsidize the transaction
	path := `m/44'/60'/3'/0/0`

	// load the ekho compiled artifact
	ekhoArtifact, err := ioutil.ReadFile("artifacts/ekho.json")
	if err != nil {
		t.Errorf("error loading ekho artifact. Error: %s", err.Error())
	}

	ekhoCompiledArtifact := nchain.CompiledArtifact{}
	err = json.Unmarshal(ekhoArtifact, &ekhoCompiledArtifact)
	if err != nil {
		t.Errorf("error converting ekho compiled artifact. Error: %s", err.Error())
	}

	tt := []struct {
		network        string
		name           string
		derivationPath string
		walletID       string
		artifact       nchain.CompiledArtifact
	}{
		{rinkebyNetworkID, "ekho", path, wallet.ID.String(), ekhoCompiledArtifact},
		//{ropstenNetworkID, "ekho", path, wallet.ID.String(), ekhoCompiledArtifact},
		//{kovanNetworkID, "ekho", path, wallet.ID.String(), ekhoCompiledArtifact},
		//{goerliNetworkID, "ekho", path, wallet.ID.String(), ekhoCompiledArtifact}, // note: goerli network has provide-go bug, likely present in bookie as well
	}

	for _, tc := range tt {

		t.Logf("creating contract using wallet id: %s, derivation path: %s", tc.walletID, tc.derivationPath)
		contract, err := nchain.CreateContract(*appToken.Token, map[string]interface{}{
			"network_id":     tc.network,
			"application_id": app.ID.String(),
			"wallet_id":      tc.walletID,
			"name":           tc.name,
			"address":        "0x",
			"params": map[string]interface{}{
				"wallet_id":          tc.walletID,
				"hd_derivation_path": tc.derivationPath,
				"compiled_artifact":  tc.artifact,
			},
		})
		if err != nil {
			t.Errorf("error creating %s contract. Error: %s", tc.name, err.Error())
			return
		}

		// wait for the contract to be deployed
		started := time.Now().Unix()
		for {
			if time.Now().Unix()-started >= contractTimeout {
				t.Errorf("timed out awaiting contract address for %s contract for %s network", tc.name, tc.network)
				return
			}

			cntrct, err := nchain.GetContractDetails(*appToken.Token, contract.ID.String(), map[string]interface{}{})
			if err != nil {
				t.Errorf("error fetching %s contract details; %s", tc.name, err.Error())
				return
			}

			if cntrct.Address != nil && *cntrct.Address != "0x" {
				t.Logf("%s contract address resolved; contract id: %s; address: %s", tc.name, cntrct.ID.String(), *cntrct.Address)
				break
			}

			t.Logf("resolving contract for %s network ...", tc.network)
			time.Sleep(contractSleepTime * time.Second)
		}

		// create a message for ekho
		msg := common.RandomString(118)
		t.Logf("executing contract using wallet id: %s, derivation path: %s", tc.walletID, tc.derivationPath)

		params := map[string]interface{}{}

		parameter := fmt.Sprintf(`{"method":"broadcast", "hd_derivation_path": "%s", "params": ["%s"], "value":0, "wallet_id":"%s"}`, tc.derivationPath, msg, tc.walletID)
		json.Unmarshal([]byte(parameter), &params)

		// execute the contract method
		execResponse, err := nchain.ExecuteContract(*appToken.Token, contract.ID.String(), params)
		if err != nil {
			t.Errorf("error executing contract. Error: %s", err.Error())
			return
		}

		if err != nil {
			t.Errorf("error executing contract: %s", err.Error())
			return
		}

		// wait for the transaction to be mined (get a tx hash)
		started = time.Now().Unix()
		for {
			if time.Now().Unix()-started >= transactionTimeout {
				t.Error("timed out awaiting transaction hash")
				return
			}

			tx, err := nchain.GetTransactionDetails(*appToken.Token, *execResponse.Reference, map[string]interface{}{})
			//this is populated by nchain consumer, so it can take a moment to appear, so we won't quit right away on a 404
			if err != nil {
				t.Logf("tx not yet mined...")
			}

			if err == nil {
				if tx.Hash != nil && *tx.Hash != "0x" {
					t.Logf("tx resolved; tx id: %s; hash: %s", tx.ID.String(), *tx.Hash)
					break
				}
				t.Logf("resolving transaction...")
			}
			time.Sleep(transactionSleepTime * time.Second)
		}

		t.Logf("contract execution successful")
	}
}

func TestContractAccount(t *testing.T) {

	t.Parallel()

	testId, err := uuid.NewV4()
	if err != nil {
		t.Logf("error creating new UUID")
	}

	userToken, err := UserAndTokenFactory(testId)
	if err != nil {
		t.Errorf("user authentication failed. Error: %s", err.Error())
	}

	testcaseApp := Application{
		"app" + testId.String(),
		"appdesc " + testId.String(),
	}

	app, err := appFactory(*userToken, testcaseApp.name, testcaseApp.description)
	if err != nil {
		t.Errorf("error setting up application. Error: %s", err.Error())
		return
	}

	appToken, err := appTokenFactory(*userToken, app.ID)
	if err != nil {
		t.Errorf("error getting app token. Error: %s", err.Error())
		return
	}

	// load the ekho compiled artifact
	ekhoArtifact, err := ioutil.ReadFile("artifacts/ekho.json")
	if err != nil {
		t.Errorf("error loading ekho artifact. Error: %s", err.Error())
	}

	ekhoCompiledArtifact := nchain.CompiledArtifact{}
	err = json.Unmarshal(ekhoArtifact, &ekhoCompiledArtifact)
	if err != nil {
		t.Errorf("error converting ekho compiled artifact. Error: %s", err.Error())
	}

	tt := []struct {
		network  string
		name     string
		artifact nchain.CompiledArtifact
	}{
		{rinkebyNetworkID, "ekho", ekhoCompiledArtifact},
		//{ropstenNetworkID, "ekho", ekhoCompiledArtifact},
		//{kovanNetworkID, "ekho", ekhoCompiledArtifact},
		//{goerliNetworkID, "ekho", ekhoCompiledArtifact}, // note: goerli network has provide-go bug, likely present in bookie as well
	}

	for _, tc := range tt {

		// create the account for that user, for the Ropsten network
		account, err := nchain.CreateAccount(*appToken.Token, map[string]interface{}{
			"network_id":     tc.network,
			"application_id": app.ID,
		})

		if err != nil {
			t.Errorf("error creating user account. Error: %s", err.Error())
		}

		t.Logf("creating contract using account id: %s", account.ID)
		contract, err := nchain.CreateContract(*appToken.Token, map[string]interface{}{
			"network_id":     tc.network,
			"application_id": app.ID.String(),
			"name":           tc.name,
			"address":        "0x",
			"params": map[string]interface{}{
				"compiled_artifact": tc.artifact,
				"account_id":        account.ID,
			},
		})
		if err != nil {
			t.Errorf("error creating %s contract. Error: %s", tc.name, err.Error())
			return
		}

		// wait for the contract to be deployed
		started := time.Now().Unix()
		for {
			if time.Now().Unix()-started >= contractTimeout {
				t.Errorf("timed out awaiting contract address for %s contract for %s network", tc.name, tc.network)
				return
			}

			cntrct, err := nchain.GetContractDetails(*appToken.Token, contract.ID.String(), map[string]interface{}{})
			if err != nil {
				t.Errorf("error fetching %s contract details; %s", tc.name, err.Error())
				return
			}

			if cntrct.Address != nil && *cntrct.Address != "0x" {
				t.Logf("%s contract address resolved; contract id: %s; address: %s", tc.name, cntrct.ID.String(), *cntrct.Address)
				break
			}

			t.Logf("resolving contract for %s network ...", tc.network)
			time.Sleep(contractSleepTime * time.Second)
		}

		// create a message for ekho
		msg := common.RandomString(118)
		t.Logf("executing contract using account id: %s", account.ID)

		params := map[string]interface{}{}

		parameter := fmt.Sprintf(`{"method":"broadcast", "params": ["%s"], "value":0, "account_id":"%s"}`, msg, account.ID)
		json.Unmarshal([]byte(parameter), &params)

		// execute the contract method
		execResponse, err := nchain.ExecuteContract(*appToken.Token, contract.ID.String(), params)
		if err != nil {
			t.Errorf("error executing contract. Error: %s", err.Error())
			return
		}

		if err != nil {
			t.Errorf("error executing contract: %s", err.Error())
			return
		}

		// wait for the transaction to be mined (get a tx hash)
		started = time.Now().Unix()
		for {
			if time.Now().Unix()-started >= transactionTimeout {
				t.Error("timed out awaiting transaction hash")
				return
			}

			tx, err := nchain.GetTransactionDetails(*appToken.Token, *execResponse.Reference, map[string]interface{}{})
			//this is populated by nchain consumer, so it can take a moment to appear, so we won't quit right away on a 404
			if err != nil {
				t.Logf("tx not yet mined...")
			}

			if err == nil {
				if tx.Hash != nil && *tx.Hash != "0x" {
					t.Logf("tx resolved; tx id: %s; hash: %s", tx.ID.String(), *tx.Hash)
					break
				}
				t.Logf("resolving transaction...")
			}
			time.Sleep(transactionSleepTime * time.Second)
		}

		t.Logf("contract execution successful")
	}
}
