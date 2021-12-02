package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/vikkoro/gocardano-api/pkg/cardano"
	"github.com/vikkoro/gocardano-api/pkg/config"
	"io"
	"math"
	"strconv"
	"strings"
)

// Service interface used to list the strings
type Service interface {
	ParsePayments(string) ([]cardano.Payment, uint64, error)
}

type service struct {
	cfg *config.Configuration
}

// NewService constructor of the default service.
func NewService(cfg *config.Configuration) *service {
	return &service{cfg}
}

// Parse SCV file into array of Payments
func (s *service) ParsePayments(cvsString string) ([]cardano.Payment, uint64, error) {

	csvReader := csv.NewReader(strings.NewReader(cvsString))

	var paymentsArray []cardano.Payment

	// We stockpile all transfers amounts here
	var totalAmountLovelace uint64

	// Here we could CSV file lines
	var lineCounter uint64 = 1

	for {
		// Read CSV file line by line
		rec, err := csvReader.Read()

		// Until the end of file
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, 0, err
		}

		// We receive amount in CSV as float. Later we will need
		// to multiply it by 1M to use in the node API
		amountAda, err := strconv.ParseFloat(rec[1], 64)
		if err != nil {
			return nil, 0, err
		}

		if amountAda < 1 {
			return nil, 0, errors.New(fmt.Sprintf("payment amount is lower than 1 ADA at line %d", lineCounter))
		}

		// Use multiplier to convert ADA to Lovelace for compatibility with the node API
		amountLovelace := uint64(math.Round(amountAda * float64(s.cfg.Multiplier)))

		// Add up transfers amounts
		totalAmountLovelace += amountLovelace

		paymentsArray = append(paymentsArray, cardano.Payment{
			Address: rec[0],
			Amount: cardano.Amount{
				Quantity: float64(amountLovelace),
				Unit:     "lovelace",
			},
		})

		// do something with read line
		//fmt.Printf("%+v\n", rec)

		lineCounter++
	}

	if lineCounter > s.cfg.PaymentsMax {
		return nil, 0, errors.New(fmt.Sprintf("number of payments is bigger, than expected: %d > %d", lineCounter, s.cfg.PaymentsMax))
	}

	return paymentsArray, totalAmountLovelace, nil
}
