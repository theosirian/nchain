## goldmine

API for building best-of-breed applications which leverage a public blockchain.

### Authentication

Consumers of this API will present a `bearer` authentication header (i.e., using a `JWT` token) for all requests. The mechanism to require this authentication has not yet been included in the codebase to simplify development and integration, but it will be included in the coming weeks; this section will be updated with specifics when authentication is required.

### Authorization

The `bearer` authorization header will be scoped to an authorized application. The `bearer` authorization header may contain a `sub` (see [RFC-7519 §§ 4.1.2](https://tools.ietf.org/html/rfc7519#section-4.1.2)) to further limit its authorized scope to a specific token or smart contract, wallet or other entity.

Certain APIs will be metered similarly to how AWS meters some of its webservices. Production applications will need a sufficient PRVD token balance to consume metered APIs (based on market conditions at the time of consumption, some quantity of PRVD tokens will be burned as a result of such metered API usage. *The PRVD token model and economics are in the process of being peer-reviewed and finalized; the documentation will be updated accordingly with specifics.*

---
The following APIs are exposed:


### Networks API

##### `GET /api/v1/networks`

Enumerate available blockchain networks and related configuration details.

```
[prvd@vpc ~]# curl -v https://goldmine.provide.services/api/v1/networks

> GET /api/v1/networks HTTP/1.1
> Host: goldmine.provide.services
> User-Agent: curl/7.51.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Thu, 28 Dec 2017 16:54:12 GMT
< Content-Type: application/json; charset=UTF-8
< Content-Length: 1899
< Connection: keep-alive
<
[
    {
        "id": "07c85f35-aa6d-4ec2-8a92-2240a85e91e9",
        "created_at": "2017-12-25T11:54:57.62033Z",
        "name": "Lightning",
        "description": "Lightning Network mainnet",
        "is_production": true,
        "sidechain_id": null,
        "config": null
    },
    {
        "id": "017428e4-41ac-41ab-bd08-8fc234e8169f",
        "created_at": "2017-12-25T11:54:57.622145Z",
        "name": "Lightning Testnet",
        "description": "Lightning Network testnet",
        "is_production": false,
        "sidechain_id": null,
        "config": null
    },
    {
        "id": "6eafb694-95a9-407f-a6d2-f541659c49b9",
        "created_at": "2017-12-25T11:54:57.615578Z",
        "name": "Bitcoin",
        "description": "Bitcoin mainnet",
        "is_production": true,
        "sidechain_id": "07c85f35-aa6d-4ec2-8a92-2240a85e91e9",
        "config": null
    },
    {
        "id": "b018af93-7c7f-4b76-a0d3-2f4282250e82",
        "created_at": "2017-12-25T11:54:57.61854Z",
        "name": "Bitcoin Testnet",
        "description": "Bitcoin testnet",
        "is_production": false,
        "sidechain_id": "017428e4-41ac-41ab-bd08-8fc234e8169f",
        "config": null
    },
    {
        "id": "5bc7d17f-653f-4599-a6dd-618ae3a1ecb2",
        "created_at": "2017-12-25T11:54:57.629505Z",
        "name": "Ethereum",
        "description": "Ethereum mainnet",
        "is_production": true,
        "sidechain_id": null,
        "config": null
    },
    {
        "id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
        "created_at": "2017-12-25T11:54:57.63379Z",
        "name": "Ethereum testnet",
        "description": "ROPSTEN (Revival) TESTNET",
        "is_production": false,
        "sidechain_id": null,
        "config": {
            "json_rpc_url": "https://ethereum-ropsten-testnet-json-rpc.provide.services",
            "testnet": "ropsten"
        }
    }
]
```

##### `GET /api/v1/networks/:id`
##### `GET /api/v1/networks/:id/addresses`
##### `GET /api/v1/networks/:id/blocks`
##### `GET /api/v1/networks/:id/contracts`
##### `GET /api/v1/networks/:id/transactions`


### Prices API

##### `GET /api/v1/prices`

Fetch real-time pricing data for major currency pairs and supported tokens.

```
[prvd@vpc ~]# curl -v https://goldmine.provide.services/api/v1/prices

> GET /api/v1/prices HTTP/1.1
> Host: goldmine.provide.services
> User-Agent: curl/7.51.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Thu, 28 Dec 2017 17:03:42 GMT
< Content-Type: application/json; charset=UTF-8
< Content-Length: 88
< Connection: keep-alive
<
{
    "btcusd": 14105.2,
    "ethusd": 707,
    "ltcusd": 244.48,
    "prvdusd": 0.22
}
```


### Contracts API

##### `GET /api/v1/contracts`

Enumerate managed smart contracts. The response contains a `params` object which includes `Network`-specific descriptors.

```
[prvd@vpc ~]# curl https://goldmine.provide.services/api/v1/contracts
[
    {
        "id": "76e6f407-735e-46e4-8281-390f770b2717",
        "created_at": "2018-01-01T19:29:38.845343Z",
        "network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
        "transaction_id": "e4a09a42-c584-4b7f-981e-c72f337b672b",
        "name": "0x0079F773651C8B7bAEAf4461C81A1494639F830F",
        "address": "0x0079F773651C8B7bAEAf4461C81A1494639F830F",
        "params": null
    },
    {
        "id": "98e1958a-aab6-45c4-ac06-45c8eaa66b57",
        "created_at": "2018-01-01T19:35:44.313535Z",
        "network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
        "transaction_id": "e050ba50-6cb2-4107-a4df-23be0782f864",
        "name": "0x7B3e4D37F4d38ec772ec0593294c269425807e8E",
        "address": "0x7B3e4D37F4d38ec772ec0593294c269425807e8E",
        "params": null
    },
    {
        "id": "4cffd49c-5cdd-4a87-8a75-ff4bd361f7ad",
        "created_at": "2018-01-01T19:55:34.916091Z",
        "network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
        "transaction_id": "a3ce7936-d35a-4b73-a961-6cd8aea509f3",
        "name": "ProvideToken",
        "address": "0x425823FA922242bBeE40BbC1aD880d5e8bb7645F",
        "params": {
            "name": "ProvideToken",
            "abi": [
                {
                    "constant": true,
                    "inputs": [
                        {
                            "name": "_holder",
                            "type": "address"
                        }
                    ],
                    "name": "tokenGrantsCount",
                    "outputs": [
                        {
                            "name": "index",
                            "type": "uint256"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "mintingFinished",
                    "outputs": [
                        {
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "name",
                    "outputs": [
                        {
                            "name": "",
                            "type": "string"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "_spender",
                            "type": "address"
                        },
                        {
                            "name": "_value",
                            "type": "uint256"
                        }
                    ],
                    "name": "approve",
                    "outputs": [],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "controller",
                            "type": "address"
                        }
                    ],
                    "name": "setUpgradeController",
                    "outputs": [],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "totalSupply",
                    "outputs": [
                        {
                            "name": "",
                            "type": "uint256"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "_from",
                            "type": "address"
                        },
                        {
                            "name": "_to",
                            "type": "address"
                        },
                        {
                            "name": "_value",
                            "type": "uint256"
                        }
                    ],
                    "name": "transferFrom",
                    "outputs": [],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [
                        {
                            "name": "",
                            "type": "address"
                        },
                        {
                            "name": "",
                            "type": "uint256"
                        }
                    ],
                    "name": "grants",
                    "outputs": [
                        {
                            "name": "granter",
                            "type": "address"
                        },
                        {
                            "name": "value",
                            "type": "uint256"
                        },
                        {
                            "name": "cliff",
                            "type": "uint64"
                        },
                        {
                            "name": "vesting",
                            "type": "uint64"
                        },
                        {
                            "name": "start",
                            "type": "uint64"
                        },
                        {
                            "name": "revokable",
                            "type": "bool"
                        },
                        {
                            "name": "burnsOnRevoke",
                            "type": "bool"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "decimals",
                    "outputs": [
                        {
                            "name": "",
                            "type": "uint256"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "_controller",
                            "type": "address"
                        }
                    ],
                    "name": "changeController",
                    "outputs": [],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "_to",
                            "type": "address"
                        },
                        {
                            "name": "_amount",
                            "type": "uint256"
                        }
                    ],
                    "name": "mint",
                    "outputs": [
                        {
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "burnAmount",
                            "type": "uint256"
                        }
                    ],
                    "name": "burn",
                    "outputs": [],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "value",
                            "type": "uint256"
                        }
                    ],
                    "name": "upgrade",
                    "outputs": [],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "upgradeAgent",
                    "outputs": [
                        {
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [
                        {
                            "name": "_holder",
                            "type": "address"
                        },
                        {
                            "name": "_grantId",
                            "type": "uint256"
                        }
                    ],
                    "name": "tokenGrant",
                    "outputs": [
                        {
                            "name": "granter",
                            "type": "address"
                        },
                        {
                            "name": "value",
                            "type": "uint256"
                        },
                        {
                            "name": "vested",
                            "type": "uint256"
                        },
                        {
                            "name": "start",
                            "type": "uint64"
                        },
                        {
                            "name": "cliff",
                            "type": "uint64"
                        },
                        {
                            "name": "vesting",
                            "type": "uint64"
                        },
                        {
                            "name": "revokable",
                            "type": "bool"
                        },
                        {
                            "name": "burnsOnRevoke",
                            "type": "bool"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [
                        {
                            "name": "holder",
                            "type": "address"
                        }
                    ],
                    "name": "lastTokenIsTransferableDate",
                    "outputs": [
                        {
                            "name": "date",
                            "type": "uint64"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [
                        {
                            "name": "_owner",
                            "type": "address"
                        }
                    ],
                    "name": "balanceOf",
                    "outputs": [
                        {
                            "name": "balance",
                            "type": "uint256"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [],
                    "name": "finishMinting",
                    "outputs": [
                        {
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "getUpgradeState",
                    "outputs": [
                        {
                            "name": "",
                            "type": "uint8"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "upgradeController",
                    "outputs": [
                        {
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "symbol",
                    "outputs": [
                        {
                            "name": "",
                            "type": "string"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "canUpgrade",
                    "outputs": [
                        {
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "_to",
                            "type": "address"
                        },
                        {
                            "name": "_value",
                            "type": "uint256"
                        },
                        {
                            "name": "_start",
                            "type": "uint64"
                        },
                        {
                            "name": "_cliff",
                            "type": "uint64"
                        },
                        {
                            "name": "_vesting",
                            "type": "uint64"
                        },
                        {
                            "name": "_revokable",
                            "type": "bool"
                        },
                        {
                            "name": "_burnsOnRevoke",
                            "type": "bool"
                        }
                    ],
                    "name": "grantVestedTokens",
                    "outputs": [],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "_to",
                            "type": "address"
                        },
                        {
                            "name": "_value",
                            "type": "uint256"
                        }
                    ],
                    "name": "transfer",
                    "outputs": [],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "totalUpgraded",
                    "outputs": [
                        {
                            "name": "",
                            "type": "uint256"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [
                        {
                            "name": "holder",
                            "type": "address"
                        },
                        {
                            "name": "time",
                            "type": "uint64"
                        }
                    ],
                    "name": "transferableTokens",
                    "outputs": [
                        {
                            "name": "",
                            "type": "uint256"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "agent",
                            "type": "address"
                        }
                    ],
                    "name": "setUpgradeAgent",
                    "outputs": [],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [
                        {
                            "name": "_owner",
                            "type": "address"
                        },
                        {
                            "name": "_spender",
                            "type": "address"
                        }
                    ],
                    "name": "allowance",
                    "outputs": [
                        {
                            "name": "remaining",
                            "type": "uint256"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [
                        {
                            "name": "tokens",
                            "type": "uint256"
                        },
                        {
                            "name": "time",
                            "type": "uint256"
                        },
                        {
                            "name": "start",
                            "type": "uint256"
                        },
                        {
                            "name": "cliff",
                            "type": "uint256"
                        },
                        {
                            "name": "vesting",
                            "type": "uint256"
                        }
                    ],
                    "name": "calculateVestedTokens",
                    "outputs": [
                        {
                            "name": "",
                            "type": "uint256"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": false,
                    "inputs": [
                        {
                            "name": "_holder",
                            "type": "address"
                        },
                        {
                            "name": "_grantId",
                            "type": "uint256"
                        }
                    ],
                    "name": "revokeTokenGrant",
                    "outputs": [],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "controller",
                    "outputs": [
                        {
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "constant": true,
                    "inputs": [],
                    "name": "BURN_ADDRESS",
                    "outputs": [
                        {
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "payable": false,
                    "type": "function"
                },
                {
                    "inputs": [],
                    "payable": false,
                    "type": "constructor"
                },
                {
                    "payable": true,
                    "type": "fallback"
                },
                {
                    "anonymous": false,
                    "inputs": [
                        {
                            "indexed": true,
                            "name": "from",
                            "type": "address"
                        },
                        {
                            "indexed": true,
                            "name": "to",
                            "type": "address"
                        },
                        {
                            "indexed": false,
                            "name": "value",
                            "type": "uint256"
                        }
                    ],
                    "name": "Upgrade",
                    "type": "event"
                },
                {
                    "anonymous": false,
                    "inputs": [
                        {
                            "indexed": false,
                            "name": "agent",
                            "type": "address"
                        }
                    ],
                    "name": "UpgradeAgentSet",
                    "type": "event"
                },
                {
                    "anonymous": false,
                    "inputs": [
                        {
                            "indexed": true,
                            "name": "from",
                            "type": "address"
                        },
                        {
                            "indexed": true,
                            "name": "to",
                            "type": "address"
                        },
                        {
                            "indexed": false,
                            "name": "value",
                            "type": "uint256"
                        },
                        {
                            "indexed": false,
                            "name": "grantId",
                            "type": "uint256"
                        }
                    ],
                    "name": "NewTokenGrant",
                    "type": "event"
                },
                {
                    "anonymous": false,
                    "inputs": [
                        {
                            "indexed": true,
                            "name": "to",
                            "type": "address"
                        },
                        {
                            "indexed": false,
                            "name": "value",
                            "type": "uint256"
                        }
                    ],
                    "name": "Mint",
                    "type": "event"
                },
                {
                    "anonymous": false,
                    "inputs": [],
                    "name": "MintFinished",
                    "type": "event"
                },
                {
                    "anonymous": false,
                    "inputs": [
                        {
                            "indexed": false,
                            "name": "burner",
                            "type": "address"
                        },
                        {
                            "indexed": false,
                            "name": "burnedAmount",
                            "type": "uint256"
                        }
                    ],
                    "name": "Burned",
                    "type": "event"
                },
                {
                    "anonymous": false,
                    "inputs": [
                        {
                            "indexed": true,
                            "name": "owner",
                            "type": "address"
                        },
                        {
                            "indexed": true,
                            "name": "spender",
                            "type": "address"
                        },
                        {
                            "indexed": false,
                            "name": "value",
                            "type": "uint256"
                        }
                    ],
                    "name": "Approval",
                    "type": "event"
                },
                {
                    "anonymous": false,
                    "inputs": [
                        {
                            "indexed": true,
                            "name": "from",
                            "type": "address"
                        },
                        {
                            "indexed": true,
                            "name": "to",
                            "type": "address"
                        },
                        {
                            "indexed": false,
                            "name": "value",
                            "type": "uint256"
                        }
                    ],
                    "name": "Transfer",
                    "type": "event"
                }
            ]
        }
    }
]
```


### Tokens API

##### `GET /api/v1/tokens`

Enumerate managed token contracts.

*This API and documentation is still being developed.*

##### `POST /api/v1/tokens`

Create a new token smart contract in accordance with the given parameters and deploy it to the specified `Network`.

*This API and documentation is still being developed.*


### Transactions API

##### `GET /api/v1/transactions`

Enumerate transactions.

*The response returned by this API will soon include network-specific metadata.*

```
[prvd@vpc ~]# curl https://goldmine.provide.services/api/v1/transactions
[
    {
        "id": "b2569500-c0d2-42bf-8992-120e7ada875d",
        "created_at": "2017-12-28T16:56:42.965056Z",
        "network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
        "wallet_id": "ce1fa3b8-049e-467b-90d8-53b9a5098b7b",
        "to": "0xfb17cB7bb99128AAb60B1DD103271d99C8237c0d",
        "value": 1000,
        "data": "",
        "hash": "a20441a6bf1f40cfc3de3238189a44af102f02aa2c97b91ae1484f7cbd9ab393"
    },
    {
        "id": "ca1e83b1-fc25-4471-8130-c53eb4e29623",
        "created_at": "2017-12-28T17:25:50.591828Z",
        "network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
        "wallet_id": "ce1fa3b8-049e-467b-90d8-53b9a5098b7b",
        "to": "0xfb17cB7bb99128AAb60B1DD103271d99C8237c0d",
        "value": 1000,
        "data": "",
        "hash": "c7b0b39276fa65801561f49adab795361c1e99e93f6f1cf727328137f7343944"
    }
]
```


##### `POST /api/v1/transactions`

Prepare and sign a protocol transaction using a managed signing `Wallet` on behalf of a specific application user and broadcast the transaction to the public blockchain `Network`. Under certain conditions, calling this API will result in a `Transaction` being created which requires lifecylce management (i.e., in the case when a managed `Sidechain` has been configured to scale micropayments channels and/or coalesce an application's transactions for on-chain settlement.

```
[prvd@vpc ~]# curl -v -XPOST -H 'content-type: application/json' https://goldmine.provide.services/api/v1/transactions \
-d '{"network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc", "wallet_id": "ce1fa3b8-049e-467b-90d8-53b9a5098b7b", "to": "0xfb17cB7bb99128AAb60B1DD103271d99C8237c0d", "value": 1000}'

> POST /api/v1/transactions HTTP/1.1
> Host: goldmine.provide.services
> User-Agent: curl/7.51.0
> Accept: */*
> Content-Length: 174
> Content-Type: application/json
>
* upload completely sent off: 174 out of 174 bytes
< HTTP/1.1 201 Created
< Date: Thu, 28 Dec 2017 16:56:43 GMT
< Content-Type: application/json; charset=UTF-8
< Content-Length: 393
< Connection: keep-alive
<
{
    "id": "b2569500-c0d2-42bf-8992-120e7ada875d",
    "created_at": "2017-12-28T16:56:42.965055765Z",
    "network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
    "wallet_id": "ce1fa3b8-049e-467b-90d8-53b9a5098b7b",
    "to": "0xfb17cB7bb99128AAb60B1DD103271d99C8237c0d",
    "value": 1000,
    "data": null,
    "hash": "a20441a6bf1f40cfc3de3238189a44af102f02aa2c97b91ae1484f7cbd9ab393"
}
```

The signed transaction is broadcast to the `Network` targeted by the given `network_id`:
![Tx broadcast to Ropsten testnet](https://s3.amazonaws.com/provide-github/ropsten-tx-example.png)

*If the broadcast transaction represents a contract deployment, a `Contract` will be created implicitly after the deployment has been confirmed with the `Network`. The following example represents a `Contract` creation with provided `params` specific to the Ethereum network:*

```
[prvd@vpc ~]# curl -v -XPOST -H 'content-type: application/json' http://goldmine.provide.services/api/v1/transactions -d '{"network_id":"ba02ff92-f5bb-4d44-9187-7e1cc214b9fc","wallet_id":"ce1fa3b8-049e-467b-90d8-53b9a5098b7b","data":"60606040526003805460a060020a60ff021916905560006004556014600555341561002957600080fd5b5b335b5b60038054600160a060020a03191633600160a060020a03161790555b60078054600160a060020a031916600160a060020a0383161790555b505b5b61204f806100776000396000f3006060604052361561017a5763ffffffff60e060020a60003504166302a72a4c811461022757806305d2035b1461025857806306fdde031461027f578063095ea7b31461030a5780630de54c081461032e57806318160ddd1461034f57806323b872dd146103745780632c71e60a1461039e578063313ce567146104195780633cebb8231461043e57806340c10f191461045f57806342966c681461049557806345977d03146104ad5780635de4ccb0146104c5578063600e85b7146104f45780636c182e991461057557806370a08231146105b15780637d64bcb4146105e25780638444b3911461060957806387543ef61461064057806395d89b411461066f5780639738968c146106fa5780639754a4d914610721578063a9059cbb14610768578063c752ff621461078c578063d347c205146107b1578063d7e7088a146107ef578063dd62ed3e14610810578063df3c211b14610847578063eb944e4c1461087b578063f77c47911461089f578063fccc2813146108ce575b6102255b60035461019390600160a060020a03166108fd565b1561021d57600354600160a060020a031663f48c3054343360006040516020015260405160e060020a63ffffffff8516028152600160a060020a0390911660048201526024016020604051808303818588803b15156101f157600080fd5b6125ee5a03f1151561020257600080fd5b5050505060405180519050151561021857600080fd5b610222565b600080fd5b5b565b005b341561023257600080fd5b610246600160a060020a036004351661092a565b60405190815260200160405180910390f35b341561026357600080fd5b61026b610949565b604051901515815260200160405180910390f35b341561028a57600080fd5b61029261096a565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156102cf5780820151818401525b6020016102b6565b50505050905090810190601f1680156102fc5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561031557600080fd5b610225600160a060020a03600435166024356109a1565b005b341561033957600080fd5b610225600160a060020a0360043516610a43565b005b341561035a57600080fd5b610246610a92565b60405190815260200160405180910390f35b341561037f57600080fd5b610225600160a060020a0360043581169060243516604435610a98565b005b34156103a957600080fd5b6103c0600160a060020a0360043516602435610ac4565b604051600160a060020a039097168752602087019590955267ffffffffffffffff93841660408088019190915292841660608701529216608085015290151560a084015290151560c083015260e0909101905180910390f35b341561042457600080fd5b610246610b4a565b60405190815260200160405180910390f35b341561044957600080fd5b610225600160a060020a0360043516610b4f565b005b341561046a57600080fd5b61026b600160a060020a0360043516602435610b8a565b604051901515815260200160405180910390f35b34156104a057600080fd5b610225600435610c6d565b005b34156104b857600080fd5b610225600435610d37565b005b34156104d057600080fd5b6104d8610ea3565b604051600160a060020a03909116815260200160405180910390f35b34156104ff57600080fd5b610516600160a060020a0360043516602435610eb2565b604051600160a060020a039098168852602088019690965260408088019590955267ffffffffffffffff9384166060880152918316608087015290911660a0850152151560c084015290151560e0830152610100909101905180910390f35b341561058057600080fd5b610594600160a060020a0360043516610fff565b60405167ffffffffffffffff909116815260200160405180910390f35b34156105bc57600080fd5b610246600160a060020a0360043516611091565b60405190815260200160405180910390f35b34156105ed57600080fd5b61026b6110b0565b604051901515815260200160405180910390f35b341561061457600080fd5b61061c611137565b6040518082600481111561062c57fe5b60ff16815260200191505060405180910390f35b341561064b57600080fd5b6104d8611188565b604051600160a060020a03909116815260200160405180910390f35b341561067a57600080fd5b610292611197565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156102cf5780820151818401525b6020016102b6565b50505050905090810190601f1680156102fc5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561070557600080fd5b61026b6111ce565b604051901515815260200160405180910390f35b341561072c57600080fd5b610225600160a060020a036004351660243567ffffffffffffffff6044358116906064358116906084351660a435151560c43515156111d4565b005b341561077357600080fd5b610225600160a060020a036004351660243561144e565b005b341561079757600080fd5b610246611478565b60405190815260200160405180910390f35b34156107bc57600080fd5b610246600160a060020a036004351667ffffffffffffffff6024351661147e565b60405190815260200160405180910390f35b34156107fa57600080fd5b610225600160a060020a03600435166115d6565b005b341561081b57600080fd5b610246600160a060020a0360043581169060243516611782565b60405190815260200160405180910390f35b341561085257600080fd5b6102466004356024356044356064356084356117af565b60405190815260200160405180910390f35b341561088657600080fd5b610225600160a060020a0360043516602435611821565b005b34156108aa57600080fd5b6104d8611c1b565b604051600160a060020a03909116815260200160405180910390f35b34156108d957600080fd5b6104d8611c2a565b604051600160a060020a03909116815260200160405180910390f35b600080600160a060020a03831615156109195760009150610924565b823b90506000811191505b50919050565b600160a060020a0381166000908152600660205260409020545b919050565b60035474010000000000000000000000000000000000000000900460ff1681565b60408051908101604052600781527f50726f7669646500000000000000000000000000000000000000000000000000602082015281565b80158015906109d45750600160a060020a0333811660009081526002602090815260408083209386168352929052205415155b156109de57600080fd5b600160a060020a03338116600081815260026020908152604080832094871680845294909152908190208490557f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259084905190815260200160405180910390a35b5050565b600160a060020a0381161515610a5857600080fd5b60075433600160a060020a03908116911614610a7357600080fd5b60078054600160a060020a031916600160a060020a0383161790555b50565b60045481565b8281610aa4824261147e565b811115610ab057600080fd5b610abb858585611c2f565b5b5b5050505050565b600660205281600052604060002081815481101515610adf57fe5b906000526020600020906003020160005b5080546001820154600290920154600160a060020a03909116935090915067ffffffffffffffff80821691680100000000000000008104821691608060020a8204169060ff60c060020a820481169160c860020a90041687565b600881565b60035433600160a060020a03908116911614610b6a57600080fd5b60038054600160a060020a031916600160a060020a0383161790555b5b50565b60035460009033600160a060020a03908116911614610ba857600080fd5b60035474010000000000000000000000000000000000000000900460ff1615610bd057600080fd5b600454610be3908363ffffffff611d4016565b600455600160a060020a038316600090815260016020526040902054610c0f908363ffffffff611d4016565b600160a060020a0384166000818152600160205260409081902092909255907f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d41213968859084905190815260200160405180910390a25060015b5b5b92915050565b33600160a060020a038116600090815260016020526040902054610c919083611d5a565b600160a060020a03821660009081526001602052604081209190915554610cbe908363ffffffff611d5a16565b6000557f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df78183604051600160a060020a03909216825260208201526040908101905180910390a16000600160a060020a0382166000805160206120048339815191528460405190815260200160405180910390a35b5050565b6000610d41611137565b905060035b816004811115610d5257fe5b1480610d6a575060045b816004811115610d6857fe5b145b1515610d7557600080fd5b811515610d8157600080fd5b600160a060020a033316600090815260016020526040902054610daa908363ffffffff611d5a16565b600160a060020a03331660009081526001602052604081209190915554610dd7908363ffffffff611d5a16565b600055600954610ded908363ffffffff611d4016565b600955600854600160a060020a031663753e88e5338460405160e060020a63ffffffff8516028152600160a060020a0390921660048301526024820152604401600060405180830381600087803b1515610e4657600080fd5b6102c65a03f11515610e5757600080fd5b5050600854600160a060020a03908116915033167f7e5c344a8141a805725cb476f76c6953b842222b967edd1f78ddb6e8b3f397ac8460405190815260200160405180910390a35b5050565b600854600160a060020a031681565b6000806000806000806000806000600660008c600160a060020a0316600160a060020a031681526020019081526020016000208a815481101515610ef257fe5b906000526020600020906003020160005b50805460018201546002830154600160a060020a039092169b50995067ffffffffffffffff608060020a820481169850808216975068010000000000000000820416955060ff60c060020a82048116955060c860020a9091041692509050610fee8160e060405190810160409081528254600160a060020a031682526001830154602083015260029092015467ffffffffffffffff8082169383019390935268010000000000000000810483166060830152608060020a8104909216608082015260ff60c060020a83048116151560a083015260c860020a909204909116151560c082015242611d71565b96505b509295985092959890939650565b600160a060020a03811660009081526006602052604081205442915b8181101561108957600160a060020a0384166000908152600660205260409020805461107e91908390811061104c57fe5b906000526020600020906003020160005b506002015468010000000000000000900467ffffffffffffffff1684611dc1565b92505b60010161101b565b5b5050919050565b600160a060020a0381166000908152600160205260409020545b919050565b60035460009033600160a060020a039081169116146110ce57600080fd5b6003805474ff00000000000000000000000000000000000000001916740100000000000000000000000000000000000000001790557fae5184fba832cb2b1f702aca6117b8d265eaf03ad33eb133f19dde0f5920fa0860405160405180910390a15060015b5b90565b60006111416111ce565b151561114f57506001611133565b600854600160a060020a0316151561116957506002611133565b600954151561117a57506003611133565b506004611133565b5b5b5b90565b600754600160a060020a031681565b60408051908101604052600481527f5052564400000000000000000000000000000000000000000000000000000000602082015281565b60015b90565b60008567ffffffffffffffff168567ffffffffffffffff16108061120b57508467ffffffffffffffff168467ffffffffffffffff16105b1561121557600080fd5b6005546112218961092a565b111561122c57600080fd5b600160a060020a03881660009081526006602052604090208054600181016112548382611f48565b916000526020600020906003020160005b60e0604051908101604052808761127d57600061127f565b335b600160a060020a03168152602081018c905267ffffffffffffffff808b16604083015289811660608301528b16608082015287151560a082015286151560c09091015291905081518154600160a060020a031916600160a060020a039190911617815560208201518160010155604082015160028201805467ffffffffffffffff191667ffffffffffffffff9290921691909117905560608201518160020160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060808201518160020160106101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060a082015160028201805491151560c060020a0278ff0000000000000000000000000000000000000000000000001990921691909117905560c08201516002909101805491151560c860020a0279ff00000000000000000000000000000000000000000000000000199092169190911790555090506113f2888861144e565b87600160a060020a031633600160a060020a03167ff9565aecd648a0466ffb964a79eeccdf1120ad6276189c687a6e9fe73984d9bb896001850360405191825260208201526040908101905180910390a35b5050505050505050565b338161145a824261147e565b81111561146657600080fd5b6114708484611df0565b5b5b50505050565b60095481565b600080600080600061148f8761092a565b93508315156114a8576114a187611091565b94506115cc565b60009250600091505b8382101561159b57600160a060020a0387166000908152600660205260409020805461158d9161158091859081106114e557fe5b906000526020600020906003020160005b5060e060405190810160409081528254600160a060020a031682526001830154602083015260029092015467ffffffffffffffff8082169383019390935268010000000000000000810483166060830152608060020a8104909216608082015260ff60c060020a83048116151560a083015260c860020a909204909116151560c082015288611eab565b849063ffffffff611d4016565b92505b6001909101906114b1565b6115b4836115a889611091565b9063ffffffff611d5a16565b90506115c9816115c48989611ed4565b611ee8565b94505b5050505092915050565b6115de6111ce565b15156115e957600080fd5b600160a060020a03811615156115fe57600080fd5b60075433600160a060020a0390811691161461161957600080fd5b60045b611624611137565b600481111561162f57fe5b141561163a57600080fd5b60088054600160a060020a031916600160a060020a038381169190911791829055166361d3d7a66000604051602001526040518163ffffffff1660e060020a028152600401602060405180830381600087803b151561169857600080fd5b6102c65a03f115156116a957600080fd5b5050506040518051905015156116be57600080fd5b600080546008549091600160a060020a0390911690634b2ba0dd90604051602001526040518163ffffffff1660e060020a028152600401602060405180830381600087803b151561170e57600080fd5b6102c65a03f1151561171f57600080fd5b5050506040518051905014151561173557600080fd5b6008547f7845d5aa74cc410e35571258d954f23b82276e160fe8c188fa80566580f279cc90600160a060020a0316604051600160a060020a03909116815260200160405180910390a15b50565b600160a060020a038083166000908152600260209081526040808320938516835292905220545b92915050565b600080838610156117c35760009150611817565b8286106117d257869150611817565b6118116117e5848763ffffffff611d5a16565b6118056117f8898963ffffffff611d5a16565b8a9063ffffffff611f0216565b9063ffffffff611f3116565b90508091505b5095945050505050565b600160a060020a03821660009081526006602052604081208054829182918590811061184957fe5b906000526020600020906003020160005b50600281015490935060c060020a900460ff16151561187857600080fd5b825433600160a060020a0390811691161461189257600080fd5b600283015460c860020a900460ff166118ab57336118ae565b60005b915061193d8360e060405190810160409081528254600160a060020a031682526001830154602083015260029092015467ffffffffffffffff8082169383019390935268010000000000000000810483166060830152608060020a8104909216608082015260ff60c060020a83048116151560a083015260c860020a909204909116151560c082015242611eab565b600160a060020a03861660009081526006602052604090208054919250908590811061196557fe5b906000526020600020906003020160005b508054600160a060020a0319168155600060018083018290556002909201805479ffffffffffffffffffffffffffffffffffffffffffffffffffff19169055600160a060020a0387168152600660205260409020805490916119de919063ffffffff611d5a16565b815481106119e857fe5b906000526020600020906003020160005b50600160a060020a0386166000908152600660205260409020805486908110611a1e57fe5b906000526020600020906003020160005b5081548154600160a060020a031916600160a060020a03918216178255600180840154908301556002928301805493909201805467ffffffffffffffff191667ffffffffffffffff94851617808255835468010000000000000000908190048616026fffffffffffffffff000000000000000019909116178082558354608060020a9081900490951690940277ffffffffffffffff000000000000000000000000000000001990941693909317808455825460ff60c060020a918290048116151590910278ff0000000000000000000000000000000000000000000000001990921691909117808555925460c860020a9081900490911615150279ff0000000000000000000000000000000000000000000000000019909216919091179091558516600090815260066020526040902080546000190190611b709082611f48565b50600160a060020a038216600090815260016020526040902054611b9a908263ffffffff611d4016565b600160a060020a038084166000908152600160205260408082209390935590871681522054611bcf908263ffffffff611d5a16565b600160a060020a038087166000818152600160205260409081902093909355908416916000805160206120048339815191529084905190815260200160405180910390a35b5050505050565b600354600160a060020a031681565b600081565b600060606064361015611c4157600080fd5b600160a060020a038086166000908152600260209081526040808320338516845282528083205493881683526001909152902054909250611c88908463ffffffff611d4016565b600160a060020a038086166000908152600160205260408082209390935590871681522054611cbd908463ffffffff611d5a16565b600160a060020a038616600090815260016020526040902055611ce6828463ffffffff611d5a16565b600160a060020a03808716600081815260026020908152604080832033861684529091529081902093909355908616916000805160206120048339815191529086905190815260200160405180910390a35b5b5050505050565b600082820183811015611d4f57fe5b8091505b5092915050565b600082821115611d6657fe5b508082035b92915050565b6000611db883602001518367ffffffffffffffff16856080015167ffffffffffffffff16866040015167ffffffffffffffff16876060015167ffffffffffffffff166117af565b90505b92915050565b60008167ffffffffffffffff168367ffffffffffffffff161015611de55781611db8565b825b90505b92915050565b60406044361015611e0057600080fd5b600160a060020a033316600090815260016020526040902054611e29908363ffffffff611d5a16565b600160a060020a033381166000908152600160205260408082209390935590851681522054611e5e908363ffffffff611d4016565b600160a060020a0380851660008181526001602052604090819020939093559133909116906000805160206120048339815191529085905190815260200160405180910390a35b5b505050565b6000611db8611eba8484611d71565b84602001519063ffffffff611d5a16565b90505b92915050565b6000611db883611091565b90505b92915050565b6000818310611de55781611db8565b825b90505b92915050565b6000828202831580611f1e5750828482811515611f1b57fe5b04145b1515611d4f57fe5b8091505b5092915050565b60008183811515611f3e57fe5b0490505b92915050565b815481835581811511611ea557600302816003028360005260206000209182019101611ea59190611fac565b5b505050565b815481835581811511611ea557600302816003028360005260206000209182019101611ea59190611fac565b5b505050565b61113391905b80821115611ffc578054600160a060020a03191681556000600182015560028101805479ffffffffffffffffffffffffffffffffffffffffffffffffffff19169055600301611fb2565b5090565b905600ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa165627a7a72305820647a44d70f10b0501747e17ced2c1a88ccccc6cfc5f4d543a7cf726fe80fb0280029", "params": {"name": "ProvideToken", "abi": [{"constant":true,"inputs":[{"name":"_holder","type":"address"}],"name":"tokenGrantsCount","outputs":[{"name":"index","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"mintingFinished","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"controller","type":"address"}],"name":"setUpgradeController","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"uint256"}],"name":"grants","outputs":[{"name":"granter","type":"address"},{"name":"value","type":"uint256"},{"name":"cliff","type":"uint64"},{"name":"vesting","type":"uint64"},{"name":"start","type":"uint64"},{"name":"revokable","type":"bool"},{"name":"burnsOnRevoke","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_controller","type":"address"}],"name":"changeController","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"mint","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"burnAmount","type":"uint256"}],"name":"burn","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"value","type":"uint256"}],"name":"upgrade","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"upgradeAgent","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"_holder","type":"address"},{"name":"_grantId","type":"uint256"}],"name":"tokenGrant","outputs":[{"name":"granter","type":"address"},{"name":"value","type":"uint256"},{"name":"vested","type":"uint256"},{"name":"start","type":"uint64"},{"name":"cliff","type":"uint64"},{"name":"vesting","type":"uint64"},{"name":"revokable","type":"bool"},{"name":"burnsOnRevoke","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"holder","type":"address"}],"name":"lastTokenIsTransferableDate","outputs":[{"name":"date","type":"uint64"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"finishMinting","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"getUpgradeState","outputs":[{"name":"","type":"uint8"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"upgradeController","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"canUpgrade","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"},{"name":"_start","type":"uint64"},{"name":"_cliff","type":"uint64"},{"name":"_vesting","type":"uint64"},{"name":"_revokable","type":"bool"},{"name":"_burnsOnRevoke","type":"bool"}],"name":"grantVestedTokens","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"totalUpgraded","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"holder","type":"address"},{"name":"time","type":"uint64"}],"name":"transferableTokens","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"agent","type":"address"}],"name":"setUpgradeAgent","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"tokens","type":"uint256"},{"name":"time","type":"uint256"},{"name":"start","type":"uint256"},{"name":"cliff","type":"uint256"},{"name":"vesting","type":"uint256"}],"name":"calculateVestedTokens","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_holder","type":"address"},{"name":"_grantId","type":"uint256"}],"name":"revokeTokenGrant","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"controller","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"BURN_ADDRESS","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"inputs":[],"payable":false,"type":"constructor"},{"payable":true,"type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Upgrade","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"agent","type":"address"}],"name":"UpgradeAgentSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"grantId","type":"uint256"}],"name":"NewTokenGrant","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Mint","type":"event"},{"anonymous":false,"inputs":[],"name":"MintFinished","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"burner","type":"address"},{"indexed":false,"name":"burnedAmount","type":"uint256"}],"name":"Burned","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]}}'

> POST /api/v1/transactions HTTP/1.1
> Host: goldmine.provide.services
> User-Agent: curl/7.51.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 24054
> Expect: 100-continue
>
< HTTP/1.1 100 Continue
* We are completely uploaded and fine
< HTTP/1.1 201 Created
< Date: Mon, 01 Jan 2018 20:20:07 GMT
< Content-Type: application/json; charset=UTF-8
< Transfer-Encoding: chunked
< Connection: keep-alive
<
{
    "id": "a1e55081-52d3-452b-bc24-fd4030317ac5",
    "created_at": "2018-01-01T20:19:39.527211009Z",
    "network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
    "wallet_id": "ce1fa3b8-049e-467b-90d8-53b9a5098b7b",
    "to": null,
    "value": 0,
    "data": "60606040526003805460a060020a60ff021916905560006004556014600555341561002957600080fd5b5b335b5b60038054600160a060020a03191633600160a060020a03161790555b60078054600160a060020a031916600160a060020a0383161790555b505b5b61204f806100776000396000f3006060604052361561017a5763ffffffff60e060020a60003504166302a72a4c811461022757806305d2035b1461025857806306fdde031461027f578063095ea7b31461030a5780630de54c081461032e57806318160ddd1461034f57806323b872dd146103745780632c71e60a1461039e578063313ce567146104195780633cebb8231461043e57806340c10f191461045f57806342966c681461049557806345977d03146104ad5780635de4ccb0146104c5578063600e85b7146104f45780636c182e991461057557806370a08231146105b15780637d64bcb4146105e25780638444b3911461060957806387543ef61461064057806395d89b411461066f5780639738968c146106fa5780639754a4d914610721578063a9059cbb14610768578063c752ff621461078c578063d347c205146107b1578063d7e7088a146107ef578063dd62ed3e14610810578063df3c211b14610847578063eb944e4c1461087b578063f77c47911461089f578063fccc2813146108ce575b6102255b60035461019390600160a060020a03166108fd565b1561021d57600354600160a060020a031663f48c3054343360006040516020015260405160e060020a63ffffffff8516028152600160a060020a0390911660048201526024016020604051808303818588803b15156101f157600080fd5b6125ee5a03f1151561020257600080fd5b5050505060405180519050151561021857600080fd5b610222565b600080fd5b5b565b005b341561023257600080fd5b610246600160a060020a036004351661092a565b60405190815260200160405180910390f35b341561026357600080fd5b61026b610949565b604051901515815260200160405180910390f35b341561028a57600080fd5b61029261096a565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156102cf5780820151818401525b6020016102b6565b50505050905090810190601f1680156102fc5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561031557600080fd5b610225600160a060020a03600435166024356109a1565b005b341561033957600080fd5b610225600160a060020a0360043516610a43565b005b341561035a57600080fd5b610246610a92565b60405190815260200160405180910390f35b341561037f57600080fd5b610225600160a060020a0360043581169060243516604435610a98565b005b34156103a957600080fd5b6103c0600160a060020a0360043516602435610ac4565b604051600160a060020a039097168752602087019590955267ffffffffffffffff93841660408088019190915292841660608701529216608085015290151560a084015290151560c083015260e0909101905180910390f35b341561042457600080fd5b610246610b4a565b60405190815260200160405180910390f35b341561044957600080fd5b610225600160a060020a0360043516610b4f565b005b341561046a57600080fd5b61026b600160a060020a0360043516602435610b8a565b604051901515815260200160405180910390f35b34156104a057600080fd5b610225600435610c6d565b005b34156104b857600080fd5b610225600435610d37565b005b34156104d057600080fd5b6104d8610ea3565b604051600160a060020a03909116815260200160405180910390f35b34156104ff57600080fd5b610516600160a060020a0360043516602435610eb2565b604051600160a060020a039098168852602088019690965260408088019590955267ffffffffffffffff9384166060880152918316608087015290911660a0850152151560c084015290151560e0830152610100909101905180910390f35b341561058057600080fd5b610594600160a060020a0360043516610fff565b60405167ffffffffffffffff909116815260200160405180910390f35b34156105bc57600080fd5b610246600160a060020a0360043516611091565b60405190815260200160405180910390f35b34156105ed57600080fd5b61026b6110b0565b604051901515815260200160405180910390f35b341561061457600080fd5b61061c611137565b6040518082600481111561062c57fe5b60ff16815260200191505060405180910390f35b341561064b57600080fd5b6104d8611188565b604051600160a060020a03909116815260200160405180910390f35b341561067a57600080fd5b610292611197565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156102cf5780820151818401525b6020016102b6565b50505050905090810190601f1680156102fc5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561070557600080fd5b61026b6111ce565b604051901515815260200160405180910390f35b341561072c57600080fd5b610225600160a060020a036004351660243567ffffffffffffffff6044358116906064358116906084351660a435151560c43515156111d4565b005b341561077357600080fd5b610225600160a060020a036004351660243561144e565b005b341561079757600080fd5b610246611478565b60405190815260200160405180910390f35b34156107bc57600080fd5b610246600160a060020a036004351667ffffffffffffffff6024351661147e565b60405190815260200160405180910390f35b34156107fa57600080fd5b610225600160a060020a03600435166115d6565b005b341561081b57600080fd5b610246600160a060020a0360043581169060243516611782565b60405190815260200160405180910390f35b341561085257600080fd5b6102466004356024356044356064356084356117af565b60405190815260200160405180910390f35b341561088657600080fd5b610225600160a060020a0360043516602435611821565b005b34156108aa57600080fd5b6104d8611c1b565b604051600160a060020a03909116815260200160405180910390f35b34156108d957600080fd5b6104d8611c2a565b604051600160a060020a03909116815260200160405180910390f35b600080600160a060020a03831615156109195760009150610924565b823b90506000811191505b50919050565b600160a060020a0381166000908152600660205260409020545b919050565b60035474010000000000000000000000000000000000000000900460ff1681565b60408051908101604052600781527f50726f7669646500000000000000000000000000000000000000000000000000602082015281565b80158015906109d45750600160a060020a0333811660009081526002602090815260408083209386168352929052205415155b156109de57600080fd5b600160a060020a03338116600081815260026020908152604080832094871680845294909152908190208490557f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259084905190815260200160405180910390a35b5050565b600160a060020a0381161515610a5857600080fd5b60075433600160a060020a03908116911614610a7357600080fd5b60078054600160a060020a031916600160a060020a0383161790555b50565b60045481565b8281610aa4824261147e565b811115610ab057600080fd5b610abb858585611c2f565b5b5b5050505050565b600660205281600052604060002081815481101515610adf57fe5b906000526020600020906003020160005b5080546001820154600290920154600160a060020a03909116935090915067ffffffffffffffff80821691680100000000000000008104821691608060020a8204169060ff60c060020a820481169160c860020a90041687565b600881565b60035433600160a060020a03908116911614610b6a57600080fd5b60038054600160a060020a031916600160a060020a0383161790555b5b50565b60035460009033600160a060020a03908116911614610ba857600080fd5b60035474010000000000000000000000000000000000000000900460ff1615610bd057600080fd5b600454610be3908363ffffffff611d4016565b600455600160a060020a038316600090815260016020526040902054610c0f908363ffffffff611d4016565b600160a060020a0384166000818152600160205260409081902092909255907f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d41213968859084905190815260200160405180910390a25060015b5b5b92915050565b33600160a060020a038116600090815260016020526040902054610c919083611d5a565b600160a060020a03821660009081526001602052604081209190915554610cbe908363ffffffff611d5a16565b6000557f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df78183604051600160a060020a03909216825260208201526040908101905180910390a16000600160a060020a0382166000805160206120048339815191528460405190815260200160405180910390a35b5050565b6000610d41611137565b905060035b816004811115610d5257fe5b1480610d6a575060045b816004811115610d6857fe5b145b1515610d7557600080fd5b811515610d8157600080fd5b600160a060020a033316600090815260016020526040902054610daa908363ffffffff611d5a16565b600160a060020a03331660009081526001602052604081209190915554610dd7908363ffffffff611d5a16565b600055600954610ded908363ffffffff611d4016565b600955600854600160a060020a031663753e88e5338460405160e060020a63ffffffff8516028152600160a060020a0390921660048301526024820152604401600060405180830381600087803b1515610e4657600080fd5b6102c65a03f11515610e5757600080fd5b5050600854600160a060020a03908116915033167f7e5c344a8141a805725cb476f76c6953b842222b967edd1f78ddb6e8b3f397ac8460405190815260200160405180910390a35b5050565b600854600160a060020a031681565b6000806000806000806000806000600660008c600160a060020a0316600160a060020a031681526020019081526020016000208a815481101515610ef257fe5b906000526020600020906003020160005b50805460018201546002830154600160a060020a039092169b50995067ffffffffffffffff608060020a820481169850808216975068010000000000000000820416955060ff60c060020a82048116955060c860020a9091041692509050610fee8160e060405190810160409081528254600160a060020a031682526001830154602083015260029092015467ffffffffffffffff8082169383019390935268010000000000000000810483166060830152608060020a8104909216608082015260ff60c060020a83048116151560a083015260c860020a909204909116151560c082015242611d71565b96505b509295985092959890939650565b600160a060020a03811660009081526006602052604081205442915b8181101561108957600160a060020a0384166000908152600660205260409020805461107e91908390811061104c57fe5b906000526020600020906003020160005b506002015468010000000000000000900467ffffffffffffffff1684611dc1565b92505b60010161101b565b5b5050919050565b600160a060020a0381166000908152600160205260409020545b919050565b60035460009033600160a060020a039081169116146110ce57600080fd5b6003805474ff00000000000000000000000000000000000000001916740100000000000000000000000000000000000000001790557fae5184fba832cb2b1f702aca6117b8d265eaf03ad33eb133f19dde0f5920fa0860405160405180910390a15060015b5b90565b60006111416111ce565b151561114f57506001611133565b600854600160a060020a0316151561116957506002611133565b600954151561117a57506003611133565b506004611133565b5b5b5b90565b600754600160a060020a031681565b60408051908101604052600481527f5052564400000000000000000000000000000000000000000000000000000000602082015281565b60015b90565b60008567ffffffffffffffff168567ffffffffffffffff16108061120b57508467ffffffffffffffff168467ffffffffffffffff16105b1561121557600080fd5b6005546112218961092a565b111561122c57600080fd5b600160a060020a03881660009081526006602052604090208054600181016112548382611f48565b916000526020600020906003020160005b60e0604051908101604052808761127d57600061127f565b335b600160a060020a03168152602081018c905267ffffffffffffffff808b16604083015289811660608301528b16608082015287151560a082015286151560c09091015291905081518154600160a060020a031916600160a060020a039190911617815560208201518160010155604082015160028201805467ffffffffffffffff191667ffffffffffffffff9290921691909117905560608201518160020160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060808201518160020160106101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060a082015160028201805491151560c060020a0278ff0000000000000000000000000000000000000000000000001990921691909117905560c08201516002909101805491151560c860020a0279ff00000000000000000000000000000000000000000000000000199092169190911790555090506113f2888861144e565b87600160a060020a031633600160a060020a03167ff9565aecd648a0466ffb964a79eeccdf1120ad6276189c687a6e9fe73984d9bb896001850360405191825260208201526040908101905180910390a35b5050505050505050565b338161145a824261147e565b81111561146657600080fd5b6114708484611df0565b5b5b50505050565b60095481565b600080600080600061148f8761092a565b93508315156114a8576114a187611091565b94506115cc565b60009250600091505b8382101561159b57600160a060020a0387166000908152600660205260409020805461158d9161158091859081106114e557fe5b906000526020600020906003020160005b5060e060405190810160409081528254600160a060020a031682526001830154602083015260029092015467ffffffffffffffff8082169383019390935268010000000000000000810483166060830152608060020a8104909216608082015260ff60c060020a83048116151560a083015260c860020a909204909116151560c082015288611eab565b849063ffffffff611d4016565b92505b6001909101906114b1565b6115b4836115a889611091565b9063ffffffff611d5a16565b90506115c9816115c48989611ed4565b611ee8565b94505b5050505092915050565b6115de6111ce565b15156115e957600080fd5b600160a060020a03811615156115fe57600080fd5b60075433600160a060020a0390811691161461161957600080fd5b60045b611624611137565b600481111561162f57fe5b141561163a57600080fd5b60088054600160a060020a031916600160a060020a038381169190911791829055166361d3d7a66000604051602001526040518163ffffffff1660e060020a028152600401602060405180830381600087803b151561169857600080fd5b6102c65a03f115156116a957600080fd5b5050506040518051905015156116be57600080fd5b600080546008549091600160a060020a0390911690634b2ba0dd90604051602001526040518163ffffffff1660e060020a028152600401602060405180830381600087803b151561170e57600080fd5b6102c65a03f1151561171f57600080fd5b5050506040518051905014151561173557600080fd5b6008547f7845d5aa74cc410e35571258d954f23b82276e160fe8c188fa80566580f279cc90600160a060020a0316604051600160a060020a03909116815260200160405180910390a15b50565b600160a060020a038083166000908152600260209081526040808320938516835292905220545b92915050565b600080838610156117c35760009150611817565b8286106117d257869150611817565b6118116117e5848763ffffffff611d5a16565b6118056117f8898963ffffffff611d5a16565b8a9063ffffffff611f0216565b9063ffffffff611f3116565b90508091505b5095945050505050565b600160a060020a03821660009081526006602052604081208054829182918590811061184957fe5b906000526020600020906003020160005b50600281015490935060c060020a900460ff16151561187857600080fd5b825433600160a060020a0390811691161461189257600080fd5b600283015460c860020a900460ff166118ab57336118ae565b60005b915061193d8360e060405190810160409081528254600160a060020a031682526001830154602083015260029092015467ffffffffffffffff8082169383019390935268010000000000000000810483166060830152608060020a8104909216608082015260ff60c060020a83048116151560a083015260c860020a909204909116151560c082015242611eab565b600160a060020a03861660009081526006602052604090208054919250908590811061196557fe5b906000526020600020906003020160005b508054600160a060020a0319168155600060018083018290556002909201805479ffffffffffffffffffffffffffffffffffffffffffffffffffff19169055600160a060020a0387168152600660205260409020805490916119de919063ffffffff611d5a16565b815481106119e857fe5b906000526020600020906003020160005b50600160a060020a0386166000908152600660205260409020805486908110611a1e57fe5b906000526020600020906003020160005b5081548154600160a060020a031916600160a060020a03918216178255600180840154908301556002928301805493909201805467ffffffffffffffff191667ffffffffffffffff94851617808255835468010000000000000000908190048616026fffffffffffffffff000000000000000019909116178082558354608060020a9081900490951690940277ffffffffffffffff000000000000000000000000000000001990941693909317808455825460ff60c060020a918290048116151590910278ff0000000000000000000000000000000000000000000000001990921691909117808555925460c860020a9081900490911615150279ff0000000000000000000000000000000000000000000000000019909216919091179091558516600090815260066020526040902080546000190190611b709082611f48565b50600160a060020a038216600090815260016020526040902054611b9a908263ffffffff611d4016565b600160a060020a038084166000908152600160205260408082209390935590871681522054611bcf908263ffffffff611d5a16565b600160a060020a038087166000818152600160205260409081902093909355908416916000805160206120048339815191529084905190815260200160405180910390a35b5050505050565b600354600160a060020a031681565b600081565b600060606064361015611c4157600080fd5b600160a060020a038086166000908152600260209081526040808320338516845282528083205493881683526001909152902054909250611c88908463ffffffff611d4016565b600160a060020a038086166000908152600160205260408082209390935590871681522054611cbd908463ffffffff611d5a16565b600160a060020a038616600090815260016020526040902055611ce6828463ffffffff611d5a16565b600160a060020a03808716600081815260026020908152604080832033861684529091529081902093909355908616916000805160206120048339815191529086905190815260200160405180910390a35b5b5050505050565b600082820183811015611d4f57fe5b8091505b5092915050565b600082821115611d6657fe5b508082035b92915050565b6000611db883602001518367ffffffffffffffff16856080015167ffffffffffffffff16866040015167ffffffffffffffff16876060015167ffffffffffffffff166117af565b90505b92915050565b60008167ffffffffffffffff168367ffffffffffffffff161015611de55781611db8565b825b90505b92915050565b60406044361015611e0057600080fd5b600160a060020a033316600090815260016020526040902054611e29908363ffffffff611d5a16565b600160a060020a033381166000908152600160205260408082209390935590851681522054611e5e908363ffffffff611d4016565b600160a060020a0380851660008181526001602052604090819020939093559133909116906000805160206120048339815191529085905190815260200160405180910390a35b5b505050565b6000611db8611eba8484611d71565b84602001519063ffffffff611d5a16565b90505b92915050565b6000611db883611091565b90505b92915050565b6000818310611de55781611db8565b825b90505b92915050565b6000828202831580611f1e5750828482811515611f1b57fe5b04145b1515611d4f57fe5b8091505b5092915050565b60008183811515611f3e57fe5b0490505b92915050565b815481835581811511611ea557600302816003028360005260206000209182019101611ea59190611fac565b5b505050565b815481835581811511611ea557600302816003028360005260206000209182019101611ea59190611fac565b5b505050565b61113391905b80821115611ffc578054600160a060020a03191681556000600182015560028101805479ffffffffffffffffffffffffffffffffffffffffffffffffffff19169055600301611fb2565b5090565b905600ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa165627a7a72305820647a44d70f10b0501747e17ced2c1a88ccccc6cfc5f4d543a7cf726fe80fb0280029",
    "hash": "6409dd7c75a3f291c266eb45e54534c68153102c4d275c9b0d37bf43994c9d2b",
    "params": {
        "name": "ProvideToken",
        "abi": [
            {
                "constant": true,
                "inputs": [
                    {
                        "name": "_holder",
                        "type": "address"
                    }
                ],
                "name": "tokenGrantsCount",
                "outputs": [
                    {
                        "name": "index",
                        "type": "uint256"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "mintingFinished",
                "outputs": [
                    {
                        "name": "",
                        "type": "bool"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "name",
                "outputs": [
                    {
                        "name": "",
                        "type": "string"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "_spender",
                        "type": "address"
                    },
                    {
                        "name": "_value",
                        "type": "uint256"
                    }
                ],
                "name": "approve",
                "outputs": [],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "controller",
                        "type": "address"
                    }
                ],
                "name": "setUpgradeController",
                "outputs": [],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "totalSupply",
                "outputs": [
                    {
                        "name": "",
                        "type": "uint256"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "_from",
                        "type": "address"
                    },
                    {
                        "name": "_to",
                        "type": "address"
                    },
                    {
                        "name": "_value",
                        "type": "uint256"
                    }
                ],
                "name": "transferFrom",
                "outputs": [],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [
                    {
                        "name": "",
                        "type": "address"
                    },
                    {
                        "name": "",
                        "type": "uint256"
                    }
                ],
                "name": "grants",
                "outputs": [
                    {
                        "name": "granter",
                        "type": "address"
                    },
                    {
                        "name": "value",
                        "type": "uint256"
                    },
                    {
                        "name": "cliff",
                        "type": "uint64"
                    },
                    {
                        "name": "vesting",
                        "type": "uint64"
                    },
                    {
                        "name": "start",
                        "type": "uint64"
                    },
                    {
                        "name": "revokable",
                        "type": "bool"
                    },
                    {
                        "name": "burnsOnRevoke",
                        "type": "bool"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "decimals",
                "outputs": [
                    {
                        "name": "",
                        "type": "uint256"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "_controller",
                        "type": "address"
                    }
                ],
                "name": "changeController",
                "outputs": [],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "_to",
                        "type": "address"
                    },
                    {
                        "name": "_amount",
                        "type": "uint256"
                    }
                ],
                "name": "mint",
                "outputs": [
                    {
                        "name": "",
                        "type": "bool"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "burnAmount",
                        "type": "uint256"
                    }
                ],
                "name": "burn",
                "outputs": [],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "value",
                        "type": "uint256"
                    }
                ],
                "name": "upgrade",
                "outputs": [],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "upgradeAgent",
                "outputs": [
                    {
                        "name": "",
                        "type": "address"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [
                    {
                        "name": "_holder",
                        "type": "address"
                    },
                    {
                        "name": "_grantId",
                        "type": "uint256"
                    }
                ],
                "name": "tokenGrant",
                "outputs": [
                    {
                        "name": "granter",
                        "type": "address"
                    },
                    {
                        "name": "value",
                        "type": "uint256"
                    },
                    {
                        "name": "vested",
                        "type": "uint256"
                    },
                    {
                        "name": "start",
                        "type": "uint64"
                    },
                    {
                        "name": "cliff",
                        "type": "uint64"
                    },
                    {
                        "name": "vesting",
                        "type": "uint64"
                    },
                    {
                        "name": "revokable",
                        "type": "bool"
                    },
                    {
                        "name": "burnsOnRevoke",
                        "type": "bool"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [
                    {
                        "name": "holder",
                        "type": "address"
                    }
                ],
                "name": "lastTokenIsTransferableDate",
                "outputs": [
                    {
                        "name": "date",
                        "type": "uint64"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [
                    {
                        "name": "_owner",
                        "type": "address"
                    }
                ],
                "name": "balanceOf",
                "outputs": [
                    {
                        "name": "balance",
                        "type": "uint256"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [],
                "name": "finishMinting",
                "outputs": [
                    {
                        "name": "",
                        "type": "bool"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "getUpgradeState",
                "outputs": [
                    {
                        "name": "",
                        "type": "uint8"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "upgradeController",
                "outputs": [
                    {
                        "name": "",
                        "type": "address"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "symbol",
                "outputs": [
                    {
                        "name": "",
                        "type": "string"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "canUpgrade",
                "outputs": [
                    {
                        "name": "",
                        "type": "bool"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "_to",
                        "type": "address"
                    },
                    {
                        "name": "_value",
                        "type": "uint256"
                    },
                    {
                        "name": "_start",
                        "type": "uint64"
                    },
                    {
                        "name": "_cliff",
                        "type": "uint64"
                    },
                    {
                        "name": "_vesting",
                        "type": "uint64"
                    },
                    {
                        "name": "_revokable",
                        "type": "bool"
                    },
                    {
                        "name": "_burnsOnRevoke",
                        "type": "bool"
                    }
                ],
                "name": "grantVestedTokens",
                "outputs": [],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "_to",
                        "type": "address"
                    },
                    {
                        "name": "_value",
                        "type": "uint256"
                    }
                ],
                "name": "transfer",
                "outputs": [],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "totalUpgraded",
                "outputs": [
                    {
                        "name": "",
                        "type": "uint256"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [
                    {
                        "name": "holder",
                        "type": "address"
                    },
                    {
                        "name": "time",
                        "type": "uint64"
                    }
                ],
                "name": "transferableTokens",
                "outputs": [
                    {
                        "name": "",
                        "type": "uint256"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "agent",
                        "type": "address"
                    }
                ],
                "name": "setUpgradeAgent",
                "outputs": [],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [
                    {
                        "name": "_owner",
                        "type": "address"
                    },
                    {
                        "name": "_spender",
                        "type": "address"
                    }
                ],
                "name": "allowance",
                "outputs": [
                    {
                        "name": "remaining",
                        "type": "uint256"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [
                    {
                        "name": "tokens",
                        "type": "uint256"
                    },
                    {
                        "name": "time",
                        "type": "uint256"
                    },
                    {
                        "name": "start",
                        "type": "uint256"
                    },
                    {
                        "name": "cliff",
                        "type": "uint256"
                    },
                    {
                        "name": "vesting",
                        "type": "uint256"
                    }
                ],
                "name": "calculateVestedTokens",
                "outputs": [
                    {
                        "name": "",
                        "type": "uint256"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": false,
                "inputs": [
                    {
                        "name": "_holder",
                        "type": "address"
                    },
                    {
                        "name": "_grantId",
                        "type": "uint256"
                    }
                ],
                "name": "revokeTokenGrant",
                "outputs": [],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "controller",
                "outputs": [
                    {
                        "name": "",
                        "type": "address"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "constant": true,
                "inputs": [],
                "name": "BURN_ADDRESS",
                "outputs": [
                    {
                        "name": "",
                        "type": "address"
                    }
                ],
                "payable": false,
                "type": "function"
            },
            {
                "inputs": [],
                "payable": false,
                "type": "constructor"
            },
            {
                "payable": true,
                "type": "fallback"
            },
            {
                "anonymous": false,
                "inputs": [
                    {
                        "indexed": true,
                        "name": "from",
                        "type": "address"
                    },
                    {
                        "indexed": true,
                        "name": "to",
                        "type": "address"
                    },
                    {
                        "indexed": false,
                        "name": "value",
                        "type": "uint256"
                    }
                ],
                "name": "Upgrade",
                "type": "event"
            },
            {
                "anonymous": false,
                "inputs": [
                    {
                        "indexed": false,
                        "name": "agent",
                        "type": "address"
                    }
                ],
                "name": "UpgradeAgentSet",
                "type": "event"
            },
            {
                "anonymous": false,
                "inputs": [
                    {
                        "indexed": true,
                        "name": "from",
                        "type": "address"
                    },
                    {
                        "indexed": true,
                        "name": "to",
                        "type": "address"
                    },
                    {
                        "indexed": false,
                        "name": "value",
                        "type": "uint256"
                    },
                    {
                        "indexed": false,
                        "name": "grantId",
                        "type": "uint256"
                    }
                ],
                "name": "NewTokenGrant",
                "type": "event"
            },
            {
                "anonymous": false,
                "inputs": [
                    {
                        "indexed": true,
                        "name": "to",
                        "type": "address"
                    },
                    {
                        "indexed": false,
                        "name": "value",
                        "type": "uint256"
                    }
                ],
                "name": "Mint",
                "type": "event"
            },
            {
                "anonymous": false,
                "inputs": [],
                "name": "MintFinished",
                "type": "event"
            },
            {
                "anonymous": false,
                "inputs": [
                    {
                        "indexed": false,
                        "name": "burner",
                        "type": "address"
                    },
                    {
                        "indexed": false,
                        "name": "burnedAmount",
                        "type": "uint256"
                    }
                ],
                "name": "Burned",
                "type": "event"
            },
            {
                "anonymous": false,
                "inputs": [
                    {
                        "indexed": true,
                        "name": "owner",
                        "type": "address"
                    },
                    {
                        "indexed": true,
                        "name": "spender",
                        "type": "address"
                    },
                    {
                        "indexed": false,
                        "name": "value",
                        "type": "uint256"
                    }
                ],
                "name": "Approval",
                "type": "event"
            },
            {
                "anonymous": false,
                "inputs": [
                    {
                        "indexed": true,
                        "name": "from",
                        "type": "address"
                    },
                    {
                        "indexed": true,
                        "name": "to",
                        "type": "address"
                    },
                    {
                        "indexed": false,
                        "name": "value",
                        "type": "uint256"
                    }
                ],
                "name": "Transfer",
                "type": "event"
            }
        ]
    }
}
```


##### `GET /api/v1/transactions/:id`


### Wallets API

##### `GET /api/v1/wallets`

Enumerate wallets used for storing cryptocurrency or tokens on behalf of users for which Provide is managing cryptographic material (i.e., for signing transactions).

```
[prvd@vpc ~]# curl -v https://goldmine.provide.services/api/v1/wallets

> GET /api/v1/wallets HTTP/1.1
> Host: goldmine.provide.services
> User-Agent: curl/7.51.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Thu, 28 Dec 2017 16:54:18 GMT
< Content-Type: application/json; charset=UTF-8
< Content-Length: 740
< Connection: keep-alive
<
[
    {
        "id": "1e3b1b0f-c756-46dd-969d-2a5da3f1d24e",
        "created_at": "2017-12-25T12:10:04.72013Z",
        "network_id": "5bc7d17f-653f-4599-a6dd-618ae3a1ecb2",
        "address": "0xEA38C255b33FB4A8aE25998842cedF865398D286"
    },
    {
        "id": "ce1fa3b8-049e-467b-90d8-53b9a5098b7b",
        "created_at": "2017-12-28T09:57:04.365995Z",
        "network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
        "address": "0x605557a2dF436B8D4f9300450A3baD7fcc3FEBf8"
    },
    {
        "id": "e61a3a6b-3873-4edc-b3f9-fa7e45b92452",
        "created_at": "2017-12-28T10:21:41.607995Z",
        "network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
        "address": "0xfb17cB7bb99128AAb60B1DD103271d99C8237c0d"
    }
]
```


##### `POST /api/v1/wallets`

Create a managed wallet capable of storing cryptocurrencies native to a specified `Network`.

```
[prvd@vpc ~]# curl -v -XPOST https://goldmine.provide.services/api/v1/wallets -d '{"network_id":"ba02ff92-f5bb-4d44-9187-7e1cc214b9fc"}'

> POST /api/v1/wallets HTTP/1.1
> Host: goldmine.provide.services
> User-Agent: curl/7.51.0
> Accept: */*
> Content-Length: 53
> Content-Type: application/json
>
* upload completely sent off: 53 out of 53 bytes
< HTTP/1.1 201 Created
< Date: Thu, 28 Dec 2017 17:36:20 GMT
< Content-Type: application/json; charset=UTF-8
< Content-Length: 224
< Connection: keep-alive
<
{
    "id": "d24bc784-32b9-4c18-9f89-110986d6a0c4",
    "created_at": "2017-12-28T17:36:20.298961785Z",
    "network_id": "ba02ff92-f5bb-4d44-9187-7e1cc214b9fc",
    "address": "0x6282e042BE5b437Bb04E800509494186c04db882"
}
```


##### `GET /api/v1/wallets/:id`

This method is not yet implemented; it will return `Network`-specific details for the requested `Wallet`.


### Status API

##### `GET /status`

The status API is used by loadbalancers to determine if the `goldmine` instance if healthy. It returns `204 No Content` when the running microservice instance is capable of handling API requests.
