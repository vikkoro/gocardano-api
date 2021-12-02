package wallet

import (
	"errors"
	"github.com/vikkoro/gocardano-api/pkg/config"
	"log"
)

// Service interface used to list the strings
type Cardano interface {
	GetWallets() ([]Wallet, error)
	GetWallet(string) (*Wallet, error)
	GetTransferFee(*Payments) (*Estimated, error)
	Transfer(*Payments) (*TransferResponse, error)
}

// Service interface used to list the strings
type Service interface {
	GetWallets() ([]Wallet, error)
	GetWallet(string) (*Wallet, error)
	//GetTransferFee(cardano.Payments) (*cardano.Estimated, error)
	Transfer([]Payment, uint64) (*TransferResponse, error)
}

type service struct {
	cfg *config.Configuration
	c   Cardano
}

// NewService constructor of the default service.
func NewService(_cfg *config.Configuration, _c Cardano) *service {
	return &service{_cfg, _c}
}

// GetWallets returns the list of the Wallets registered.
func (s *service) GetWallets() ([]Wallet, error) {
	sl, err := s.c.GetWallets()
	if err != nil {
		log.Printf("WALLET ERROR: %q", err.Error())
		return nil, err
	}

	return sl, nil
}

// GetWallet returns an Wallet filtered by ID if the Wallet doesn't exist returns null.
func (s *service) GetWallet(id string) (*Wallet, error) {
	st, err := s.c.GetWallet(id)
	if err != nil {
		log.Printf("WALLET ERROR: %q", err.Error())
		return nil, err
	}

	return st, nil
}

// Transfer funds collected in Payment array
func (s *service) Transfer(pp []Payment, totalAmount uint64) (*TransferResponse, error) {

	// Structure to check fees and send transfers
	payments := &Payments{
		Passphrase: s.cfg.Passphrase,
		Payments:   pp,
		TimeToLive: Amount{
			Quantity: 500,
			Unit:     "second",
		},
	}

	// Check our wallet current amount
	w, err := s.c.GetWallet(s.cfg.WalletId)
	if err != nil {
		log.Printf("WALLET ERROR: %q", err.Error())
		return nil, err
	}

	// Get estimated fees for the transfers
	estimated, err := s.c.GetTransferFee(payments)
	if err != nil {
		log.Printf("WALLET ERROR: %q", err.Error())
		return nil, err
	}

	// Check if wallet amount is enough
	if w.Balance.Available.Quantity < estimated.EstimatedMax.Quantity+float64(totalAmount) {
		return nil, errors.New("not enough funds")
	}

	response, err := s.c.Transfer(payments)
	if err != nil {
		log.Printf("WALLET ERROR: %q", err.Error())
		return nil, err
	}

	return response, nil
}
