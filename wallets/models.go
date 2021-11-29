package wallets

// Error data structure
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

type CreateWalletData struct {
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

type WalletData struct {
	Name    string      `json:"name"`
	Id      string      `json:"id"`
	State   StateData   `json:"state"`
	Balance BalanceData `json:"balance"`
}

type StateData struct {
	Status string `json:"status"`
}

type BalanceData struct {
	Reward    AmountData `json:"reward"`
	Total     AmountData `json:"total"`
	Available AmountData `json:"available"`
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

type BulkPaymentsData struct {
	Passphrase string        `json:"passphrase"`
	Payments   []PaymentData `json:"payments"`
	TimeToLive AmountData    `json:"time_to_live"`
}

type PaymentData struct {
	Address string     `json:"address"`
	Amount  AmountData `json:"amount"`
}

type AmountData struct {
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
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

type TransferFundsResponseData struct {
	Id             string           `json:"id"`
	Status         string           `json:"status"`
	Direction      string           `json:"direction"`
	Amount         AmountData       `json:"amount"`
	Fee            AmountData       `json:"fee"`
	Deposit        AmountData       `json:"deposit"`
	Inputs         []InputData      `json:"inputs"`
	Outputs        []OutputData     `json:"outputs"`
	ExpiresAt      ExpiresAtData    `json:"expires_at"`
	PendingSince   PendingSinceData `json:"pending_since"`
	ScriptValidity string           `json:"script_validity"`
}

type InputData struct {
	Id      string     `json:"id"`
	Index   uint64     `json:"index"`
	Address string     `json:"address"`
	Amount  AmountData `json:"amount"`
}

type OutputData struct {
	Address string     `json:"address"`
	Amount  AmountData `json:"amount"`
}

type ExpiresAtData struct {
	Time               string `json:"time"`
	EpochNumber        uint64 `json:"epoch_number"`
	AbsoluteSlotNumber uint64 `json:"absolute_slot_number"`
	SlotNumber         uint64 `json:"slot_number"`
}

type PendingSinceData struct {
	Height             AmountData `json:"height"`
	Time               string     `json:"time"`
	EpochNumber        uint64     `json:"epoch_number"`
	AbsoluteSlotNumber uint64     `json:"absolute_slot_number"`
	SlotNumber         uint64     `json:"slot_number"`
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

type EstimatedData struct {
	EstimatedMin AmountData   `json:"estimated_min"`
	EstimatedMax AmountData   `json:"estimated_max"`
	Deposit      AmountData   `json:"deposit"`
	MinimumCoins []AmountData `json:"minimum_coins"`
}
