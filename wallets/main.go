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

func GetWallet(url string, walletId string) (*WalletData, error) {
	response, err := http.Get(url + "/" + walletId)
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

	var wallet *WalletData

	err = json.Unmarshal(body, &wallet)
	if err != nil {
		return nil, err
	}

	return wallet, err
}

func GetTransferFee(url string, walletId string, payments BulkPaymentsData) (*EstimatedData, error) {
	requestBody, err := json.Marshal(payments)

	if err != nil {
		return nil, err
	}

	response, err := http.Post(url+"/"+walletId+"/payment-fees", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var est *EstimatedData

	err = json.Unmarshal(body, &est)
	if err != nil {
		return nil, err
	}

	return est, nil
}

func TransferFunds(url string, walletId string, payments BulkPaymentsData) (*TransferFundsResponseData, error) {
	requestBody, err := json.Marshal(payments)

	if err != nil {
		return nil, err
	}

	response, err := http.Post(url+"/"+walletId+"/transactions", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var resp *TransferFundsResponseData

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
