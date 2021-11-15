package handlers

type Error struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}

type ClientHandler struct {
	Configuration
}

type Configuration struct {
	WalletId        string   `json:"wallet_id"`
	Mnemonic        []string `json:"mnemonic"`
	Passphrase      string   `json:"passphrase"`
	InformationUrl  string   `json:"information_url"`
	WalletsUrl      string   `json:"wallets_url"`
	Multiplier      uint64   `json:"multiplier"`
	UploadDirectory string   `json:"upload_directory"`
}

type ViewData struct {
	Title string
	Data  string
}
