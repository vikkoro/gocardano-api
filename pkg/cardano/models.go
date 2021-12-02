package cardano

type Error struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}

/* Data structure for the Create Wallet Request

{
    "name": "Test Wallet",
    "mnemonic_sentence": ["have", "good", "day", ... ],
    "passphrase": "1234567890",
    "address_pool_gap": 20
}
*/

type CreateWallet struct {
	Name             string   `json:"name"`
	MnemonicSentence []string `json:"mnemonic_sentence"`
	Passphrase       string   `json:"passphrase"`
	AddressPoolGap   uint64   `json:"address_pool_gap"`
}

/* Data structure for the Get Wallet Response
URL http://localhost:8090/v2/wallets/WalletId00000000000000000000000000000000

{
    "passphrase": {
        "last_updated_at": "2021-11-09T20:42:37.1684525Z"
    },
    "address_pool_gap": 20,
    "state": {
        "status": "ready"
    },
    "balance": {
        "reward": {
            "quantity": 0,
            "unit": "lovelace"
        },
        "total": {
            "quantity": 1000000000,
            "unit": "lovelace"
        },
        "available": {
            "quantity": 1000000000,
            "unit": "lovelace"
        }
    },
    "name": "Test Wallet",
    "delegation": {
        "next": [],
        "active": {
            "status": "not_delegating"
        }
    },
    "id": "WalletId00000000000000000000000000000000",
    "tip": {
        "height": {
            "quantity": 3064233,
            "unit": "block"
        },
        "time": "2021-11-11T12:13:12Z",
        "epoch_number": 168,
        "absolute_slot_number": 42263576,
        "slot_number": 57176
    },
    "assets": {
        "total": [],
        "available": []
    }
}

*/

type Wallet struct {
	Name    string  `json:"name"`
	Id      string  `json:"id"`
	State   State   `json:"state"`
	Balance Balance `json:"balance"`
}

type State struct {
	Status string `json:"status"`
}

type Balance struct {
	Reward    Amount `json:"reward"`
	Total     Amount `json:"total"`
	Available Amount `json:"available"`
}

type Amount struct {
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

/* Data structure for the Send Transfer Request
URL http://localhost:8090/v2/wallets/WalletId00000000000000000000000000000000/transactions

{
    "passphrase":"<passphrase>",
    "payments": [
      {
        "address": "addr_test000000000000000000000000000000000000000000000000000000",
        "amount": {
          "quantity": 1000000,
          "unit": "lovelace"
        }
      }
    ],
    "metadata": {
        "0": {
            "map": [
                {
                    "k": {
                        "string": "Reason"
                    },
                    "v": {
                        "string": "Have a nice day"
                    }
                }
            ]
        }
    },
    "time_to_live": {
      "quantity": 500,
      "unit": "second"
    }
  }
*/

type Payments struct {
	Passphrase string    `json:"passphrase"`
	Payments   []Payment `json:"payments"`
	TimeToLive Amount    `json:"time_to_live"`
}

type Payment struct {
	Address string `json:"address"`
	Amount  Amount `json:"amount"`
}

/* Data structure for the Payment Fees Response
URL http://localhost:8090/v2/wallets/WalletId000000000000000000000000000000/payment-fees

{
    "estimated_min": {
        "quantity": 167789,
        "unit": "lovelace"
    },
    "deposit": {
        "quantity": 0,
        "unit": "lovelace"
    },
    "minimum_coins": [
        {
            "quantity": 999978,
            "unit": "lovelace"
        }
    ],
    "estimated_max": {
        "quantity": 167789,
        "unit": "lovelace"
    }
}
*/

type Estimated struct {
	EstimatedMin Amount   `json:"estimated_min"`
	EstimatedMax Amount   `json:"estimated_max"`
	Deposit      Amount   `json:"deposit"`
	MinimumCoins []Amount `json:"minimum_coins"`
}

/* Data structure for the Send Transaction Response
URL http://localhost:8090/v2/wallets/WalletId00000000000000000000000000000000/transactions

{
    "status": "pending",
    "withdrawals": [],
    "amount": {
        "quantity": 1170473,
        "unit": "lovelace"
    },
    "inputs": [
        {
            "amount": {
                "quantity": 1000000000,
                "unit": "lovelace"
            },
            "address": "addr_test000000000000000000000000000000000000000000000000000000",
            "id": "id00000000000000000000000000000000000000000000000000000000000000",
            "assets": [],
            "index": 0
        }
    ],
    "direction": "outgoing",
    "fee": {
        "quantity": 170473,
        "unit": "lovelace"
    },
    "outputs": [
        {
            "amount": {
                "quantity": 1000000,
                "unit": "lovelace"
            },
            "address": "addr_test000000000000000000000000000000000000000000000000000000",
            "assets": []
        },
        {
            "amount": {
                "quantity": 998829527,
                "unit": "lovelace"
            },
            "address": "addr_test000000000000000000000000000000000000000000000000000000",
            "assets": []
        }
    ],
    "script_validity": "valid",
    "expires_at": {
        "time": "2021-11-11T14:08:41Z",
        "epoch_number": 168,
        "absolute_slot_number": 42270505,
        "slot_number": 64105
    },
    "pending_since": {
        "height": {
            "quantity": 3064435,
            "unit": "block"
        },
        "time": "2021-11-11T14:00:17Z",
        "epoch_number": 168,
        "absolute_slot_number": 42270001,
        "slot_number": 63601
    },
    "metadata": {
        "0": {
            "map": [
                {
                    "k": {
                        "string": "Reason"
                    },
                    "v": {
                        "string": "Have a nice day"
                    }
                }
            ]
        }
    },
    "id": "id00000000000000000000000000000000000000000000000000000000000000",
    "deposit": {
        "quantity": 0,
        "unit": "lovelace"
    },
    "collateral": [],
    "mint": []
}
*/

type TransferResponse struct {
	Id             string       `json:"id"`
	Status         string       `json:"status"`
	Direction      string       `json:"direction"`
	Amount         Amount       `json:"amount"`
	Fee            Amount       `json:"fee"`
	Deposit        Amount       `json:"deposit"`
	Inputs         []Input      `json:"inputs"`
	Outputs        []Output     `json:"outputs"`
	ExpiresAt      ExpiresAt    `json:"expires_at"`
	PendingSince   PendingSince `json:"pending_since"`
	ScriptValidity string       `json:"script_validity"`
}

type Input struct {
	Id      string `json:"id"`
	Index   uint64 `json:"index"`
	Address string `json:"address"`
	Amount  Amount `json:"amount"`
}

type Output struct {
	Address string `json:"address"`
	Amount  Amount `json:"amount"`
}

type ExpiresAt struct {
	Time               string `json:"time"`
	EpochNumber        uint64 `json:"epoch_number"`
	AbsoluteSlotNumber uint64 `json:"absolute_slot_number"`
	SlotNumber         uint64 `json:"slot_number"`
}

type PendingSince struct {
	Height             Amount `json:"height"`
	Time               string `json:"time"`
	EpochNumber        uint64 `json:"epoch_number"`
	AbsoluteSlotNumber uint64 `json:"absolute_slot_number"`
	SlotNumber         uint64 `json:"slot_number"`
}
