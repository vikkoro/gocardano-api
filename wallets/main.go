package wallets

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetList(url string) []WalletData {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	var wallets []WalletData

	err = json.Unmarshal(body, &wallets)
	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	return wallets
}

func GetWallet(url string, walletId string) WalletData {
	response, err := http.Get(url + "/" + walletId)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	var w WalletData

	err = json.Unmarshal(body, &w)
	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	return w
}

func GetTransactionFee(url string, walletId string, transactions TransactionsData) EstimatedData {
	requestBody, err := json.Marshal(transactions)

	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	response, err := http.Post(url+"/"+walletId+"/payment-fees", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	var est EstimatedData

	err = json.Unmarshal(body, &est)
	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	return est
}

func SendTransaction(url string, walletId string, transactions TransactionsData) SendTransactionsResponseData {
	requestBody, err := json.Marshal(transactions)

	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	response, err := http.Post(url+"/"+walletId+"/transactions", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	var resp SendTransactionsResponseData

	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatalf("err was %v\n", err)
	}

	return resp
}
