package parser

import (
	"github.com/vikkoro/gocardano-api/pkg/config"
	"testing"
)

var cfg *config.Configuration

var csvString = "addr_test1vz3vh7nagum5lf66ej873wur740qm32ek536gqa2wl0n24crf4ry5,2.000000\naddr_test1vqn70eqljagvahj5cmnrprrywuflu5u9k6zum75c9qdqvgcz44msn,1.500000\naddr_test1vpcsfa78jy5qwr40kzpr2s3ky7ga68al0qcsxgd9rkwnrecjvehy7,2.500000"

func setUp() {
	cfg = config.NewConfig("../../conf.json", "../../.env")
}

func tearDown() {

}

func TestMain(m *testing.M) {
	setUp()
	m.Run()
	tearDown()
}

func TestService_ParsePayments(t *testing.T) {
	// Given
	s := NewService(cfg)

	// When
	payments, amount, err := s.ParsePayments(csvString)

	// Then
	if err != nil {
		t.Errorf("Not error expected, but got: %q", err.Error())
	}

	if err == nil && len(payments) != 3 {
		t.Errorf("Did not get expected result. Wanted to read number of payments: %d, got: %d", 3, len(payments))
	}

	if err == nil && payments[2].Amount.Quantity != 2500000 {
		t.Errorf("Did not get expected result. Wanted to get amount: %d, got: %f", 2500000, payments[2].Amount.Quantity)
	}

	if err == nil && amount != 6000000 {
		t.Errorf("Did not get expected result. Wanted to get total amount: %d, got: %d", 6000000, amount)
	}

}
