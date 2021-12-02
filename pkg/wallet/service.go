package wallet

import (
	"errors"
	"github.com/vikkoro/gocardano-api/pkg/cardano"
	"github.com/vikkoro/gocardano-api/pkg/config"
	"log"
)

// Service interface used to list the strings
type Service interface {
	GetWallets() ([]cardano.Wallet, error)
	GetWallet(string) (*cardano.Wallet, error)
	//GetTransferFee(cardano.Payments) (*cardano.Estimated, error)
	Transfer([]cardano.Payment, uint64) (*cardano.TransferResponse, error)
}

type service struct {
	cfg *config.Configuration
	c   cardano.Service
}

// NewService constructor of the default service.
func NewService(_cfg *config.Configuration, _c cardano.Service) *service {
	return &service{_cfg, _c}
}

// GetWallets returns the list of the Wallets registered.
func (s *service) GetWallets() ([]cardano.Wallet, error) {
	sl, err := s.c.GetWallets()
	if err != nil {
		log.Printf("WALLET ERROR: %q", err.Error())
		return nil, err
	}

	return sl, nil
}

// GetWallet returns an Wallet filtered by ID if the Wallet doesn't exist returns null.
func (s *service) GetWallet(id string) (*cardano.Wallet, error) {
	st, err := s.c.GetWallet(id)
	if err != nil {
		log.Printf("WALLET ERROR: %q", err.Error())
		return nil, err
	}

	return st, nil
}

// Transfer funds collected in Payment array
func (s *service) Transfer(pp []cardano.Payment, totalAmount uint64) (*cardano.TransferResponse, error) {

	// Structure to check fees and send transfers
	payments := &cardano.Payments{
		Passphrase: s.cfg.Passphrase,
		Payments:   pp,
		TimeToLive: cardano.Amount{
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
