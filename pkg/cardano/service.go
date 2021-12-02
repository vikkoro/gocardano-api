package cardano

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/vikkoro/gocardano-api/pkg/config"
	"github.com/vikkoro/gocardano-api/pkg/wallet"
	"io/ioutil"
	"net/http"
)

// Service interface used to list the strings
type Service interface {
	GetWallets() ([]wallet.Wallet, error)
	GetWallet(string) (*wallet.Wallet, error)
	GetTransferFee(*wallet.Payments) (*wallet.Estimated, error)
	Transfer(*wallet.Payments) (*wallet.TransferResponse, error)
}

type service struct {
	cfg *config.Configuration
}

// NewService constructor of the default service.
func NewService(c *config.Configuration) *service {
	return &service{c}
}

func (c *service) GetWallets() ([]wallet.Wallet, error) {
	response, err := http.Get(c.cfg.WalletsUrl)
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

	if response.StatusCode != http.StatusOK {
		var errorStruct *wallet.Error
		_ = json.Unmarshal(body, &errorStruct)
		return nil, errors.New(errorStruct.Message)
	}

	var wallets []wallet.Wallet

	err = json.Unmarshal(body, &wallets)
	if err != nil {
		return nil, err
	}

	return wallets, nil

}

func (c *service) GetWallet(walletId string) (*wallet.Wallet, error) {

	response, err := http.Get(c.cfg.WalletsUrl + "/" + walletId)
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

	if response.StatusCode != http.StatusOK {
		var errorStruct *wallet.Error
		_ = json.Unmarshal(body, &errorStruct)
		return nil, errors.New(errorStruct.Message)
	}

	var w *wallet.Wallet

	err = json.Unmarshal(body, &w)
	if err != nil {
		return nil, err
	}

	return w, err

}

func (c *service) GetTransferFee(payments *wallet.Payments) (*wallet.Estimated, error) {
	requestBody, err := json.Marshal(payments)

	if err != nil {
		return nil, err
	}

	response, err := http.Post(c.cfg.WalletsUrl+"/"+c.cfg.WalletId+"/payment-fees", "application/json", bytes.NewBuffer(requestBody))

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

	if !(response.StatusCode == http.StatusOK || response.StatusCode == http.StatusAccepted) {
		var errorStruct *wallet.Error
		_ = json.Unmarshal(body, &errorStruct)
		return nil, errors.New(errorStruct.Message)
	}

	var est *wallet.Estimated

	err = json.Unmarshal(body, &est)
	if err != nil {
		return nil, err
	}

	return est, nil
}

func (c *service) Transfer(payments *wallet.Payments) (*wallet.TransferResponse, error) {
	requestBody, err := json.Marshal(payments)

	if err != nil {
		return nil, err
	}

	response, err := http.Post(c.cfg.WalletsUrl+"/"+c.cfg.WalletId+"/transactions", "application/json", bytes.NewBuffer(requestBody))

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

	if !(response.StatusCode == http.StatusOK || response.StatusCode == http.StatusAccepted) {
		var errorStruct *wallet.Error
		_ = json.Unmarshal(body, &errorStruct)
		return nil, errors.New(errorStruct.Message)
	}

	var resp *wallet.TransferResponse

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
