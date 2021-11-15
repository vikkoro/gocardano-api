# Simple REST API to interact with Cardano blockchain using Golang


## Tools used
Postman: https://www.postman.com/downloads/

You can import the `Cardano Wallet.postman_collection.json` file to see used endpoints  

## Instructions

### Download Binaries
https://github.com/input-output-hk/cardano-wallet/releases


### Download Configs
https://hydra.iohk.io/build/5047936/download/1/index.html

### Start cardano-node
run `./bin/cardano-node.exe run --port 6000 --database-path ./db-testnet --socket-path \\.\pipe\cardano-node-testnet --config ./configs/testnet-config.json --topology ./configs/testnet-topology.json`

### Serve the cardano-wallet api
run `./bin/cardano-wallet.exe serve --testnet ./configs/testnet-byron-genesis.json --node-socket \\.\pipe\cardano-node-testnet --database ./wallets-testnet --log-level INFO`

### Create a wallet
Docs: https://input-output-hk.github.io/cardano-wallet/api/edge/#tag/Wallets

run: `./bin/cardano-address.exe recovery-phrase generate`

Use the Wallet ID in `./conf.json` file

### Run the app

```bash
go run main.go
```

### Upload CSV file with addresses and amounts

Go to http://localhost:8080/api/v1/cardano/home to upload transactions CSV file in this format

		addr_test100000000000000000000000000000000000000000000000000001,1.956444
		addr_test100000000000000000000000000000000000000000000000000002,1.845180
		addr_test100000000000000000000000000000000000000000000000000003,2.395366

