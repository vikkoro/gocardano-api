package wallets

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetList(url string) ([]WalletData, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var wallets []WalletData

	err = json.Unmarshal(body, &wallets)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func GetWallet(url string, walletId string) (WalletData, error) {
	response, err := http.Get(url + "/" + walletId)
	if err != nil {
		return WalletData{}, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return WalletData{}, err
	}

	var w WalletData

	err = json.Unmarshal(body, &w)
	if err != nil {
		return WalletData{}, err
	}

	return w, err
}

func GetTransactionFee(url string, walletId string, transactions TransactionsData) (EstimatedData, error) {
	requestBody, err := json.Marshal(transactions)

	if err != nil {
		return EstimatedData{}, err
	}

	response, err := http.Post(url+"/"+walletId+"/payment-fees", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return EstimatedData{}, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return EstimatedData{}, err
	}

	var est EstimatedData

	err = json.Unmarshal(body, &est)
	if err != nil {
		return EstimatedData{}, err
	}

	return est, nil
}

func SendTransaction(url string, walletId string, transactions TransactionsData) (SendTransactionsResponseData, error) {
	requestBody, err := json.Marshal(transactions)

	if err != nil {
		return SendTransactionsResponseData{}, err
	}

	response, err := http.Post(url+"/"+walletId+"/transactions", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return SendTransactionsResponseData{}, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return SendTransactionsResponseData{}, err
	}

	var resp SendTransactionsResponseData

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return SendTransactionsResponseData{}, err
	}

	return resp, nil
}
