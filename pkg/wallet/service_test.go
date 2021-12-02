package wallet

import (
	"github.com/vikkoro/gocardano-api/pkg/cardano"
	"github.com/vikkoro/gocardano-api/pkg/config"
	"testing"
)

var cfg *config.Configuration

var payments []Payment

func setUp() {
	cfg = config.NewConfig("../../conf.json", "../../.env")

	payments = append(payments, Payment{
		Address: "addr_test1vz3vh7nagum5lf66ej873wur740qm32ek536gqa2wl0n24crf4ry5",
		Amount: Amount{
			Quantity: 2000000,
			Unit:     "lovelace",
		},
	})

	payments = append(payments, Payment{
		Address: "addr_test1vqn70eqljagvahj5cmnrprrywuflu5u9k6zum75c9qdqvgcz44msn",
		Amount: Amount{
			Quantity: 1500000,
			Unit:     "lovelace",
		},
	})

	payments = append(payments, Payment{
		Address: "addr_test1vpcsfa78jy5qwr40kzpr2s3ky7ga68al0qcsxgd9rkwnrecjvehy7",
		Amount: Amount{
			Quantity: 2500000,
			Unit:     "lovelace",
		},
	})

}

func tearDown() {

}

func TestMain(m *testing.M) {
	setUp()
	m.Run()
	tearDown()
}

func TestService_GetWallets(t *testing.T) {
	// Given
	cs := cardano.NewService(cfg)
	s := NewService(cfg, cs)

	// When
	obtained, err := s.GetWallets()

	// Then
	if err != nil {
		t.Errorf("Not error expected, but got: %q", err.Error())
	}

	if len(obtained) < 1 {
		t.Errorf("Did not get expected result. Wanted at least: 1, got: #{len(obtained)}")
	}

	if obtained[0].Id != cfg.WalletId {
		t.Errorf("Did not get expected result. Wanted: %q, got: %q", cfg.WalletId, obtained[0].Id)
	}
}

func TestService_GetWallet(t *testing.T) {
	// Given
	cs := cardano.NewService(cfg)
	s := NewService(cfg, cs)

	// When
	obtained, err := s.GetWallet(cfg.WalletId)

	// Then
	if err != nil {
		t.Errorf("Not error expected, but got: %q", err.Error())
	}

	if obtained == nil {
		t.Errorf("Did not get expected result. Wanted the get one wallet , got: nil")
	}

	if obtained != nil && obtained.Id != cfg.WalletId {
		t.Errorf("Did not get expected result. Wanted: %q, got: %q", cfg.WalletId, obtained.Id)
	}
}

func TestService_Transfer(t *testing.T) {
	// Given
	cs := cardano.NewService(cfg)
	s := NewService(cfg, cs)

	// When
	obtained, err := s.Transfer(payments, 6000000)

	// Then
	if err != nil {
		t.Errorf("Not error expected, but got: %q", err.Error())
	}

	if obtained == nil {
		t.Errorf("Did not get expected result. Wanted the get a transfer estimation result, got: nil")
	}

	if obtained != nil && !(obtained.Amount.Quantity > 6170000 && obtained.Amount.Quantity < 6210000) {
		t.Errorf("Did not get expected result. Wanted Amount be between: %d and %d, got: %f",
			6170000, 6210000, obtained.Amount.Quantity)
	}

	if obtained != nil && !(obtained.Fee.Quantity > 170000 && obtained.Fee.Quantity < 210000) {
		t.Errorf("Did not get expected result. Wanted Amount be between: %d and %d, got: %f",
			170000, 210000, obtained.Fee.Quantity)
	}
}
