package config

type Configuration struct {
	WalletId        string   `json:"wallet_id"`
	Mnemonic        []string `json:"mnemonic"`
	Passphrase      string   `json:"passphrase"`
	InformationUrl  string   `json:"information_url"`
	WalletsUrl      string   `json:"wallets_url"`
	Multiplier      uint64   `json:"multiplier"`
	UploadDirectory string   `json:"upload_directory"`
	PaymentsMax     uint64   `json:"payments_max"`
}
