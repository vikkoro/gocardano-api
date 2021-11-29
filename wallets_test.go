package main

import (
	"encoding/csv"
	"fmt"
	Handlers "github.com/vikkoro/gocardano-api/handler"
	Wallets "github.com/vikkoro/gocardano-api/wallets"
	"log"
	"os"
	"testing"
)

func TestGenerateAddress(t *testing.T) {
	configuration, err := GetConfig("conf.json")
	if err != nil {
		log.Fatal(err)
	}

	file, _ := os.Open("./data/test_small_2.csv")

	defer func() {
		_ = file.Close()
	}()

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(file)

	// Parse SCV file into array of Payments
	pp, totalAmount, err := Handlers.ParseCSVFile(csvReader, configuration.Multiplier)
	if err != nil {
		log.Fatal(err)
	}

	// Structure to check fees and send transfer
	payments := &Wallets.BulkPaymentsData{
		Passphrase: configuration.Passphrase,
		Payments:   pp,
		TimeToLive: Wallets.AmountData{
			Quantity: 500,
			Unit:     "second",
		},
	}

	// Check our wallet current amount
	wallet, err := Wallets.GetWallet(configuration.WalletsUrl, configuration.WalletId)

	if err != nil {
		log.Fatal(err)
	}

	balanceBefore := wallet.Balance.Available.Quantity

	// Get estimated fees for the transfers
	estimated, err := Wallets.GetTransferFee(configuration.WalletsUrl, configuration.WalletId, *payments)

	if err != nil {
		log.Fatal(err)
	}

	// Check if wallet amount is enough
	if wallet.Balance.Available.Quantity < estimated.EstimatedMax.Quantity+totalAmount {
		log.Fatal("Not enough funds")
	}

	// Send transactions
	_, err = Wallets.TransferFunds(configuration.WalletsUrl, configuration.WalletId, *payments)

	wallet, err = Wallets.GetWallet(configuration.WalletsUrl, configuration.WalletId)

	if err != nil {
		log.Fatal(err)
	}

	balanceAfter := wallet.Balance.Available.Quantity

	fmt.Printf("Amount %f \n", totalAmount)
	fmt.Printf("Before %f \n", balanceBefore)
	fmt.Printf("After %f \n", balanceAfter)
	fmt.Printf("Fee min %f \n", estimated.EstimatedMin.Quantity)
	fmt.Printf("Fee max %f \n", estimated.EstimatedMax.Quantity)

	//t.Errorf("Sum was incorrect, got: %d, want: %d.", 33, 10)

}

func TestParseCSVFile(t *testing.T) {

	configuration, err := GetConfig("conf.json")
	if err != nil {
		log.Fatal(err)
	}

	file, _ := os.Open("./data/test_small_2.csv")

	defer func() {
		_ = file.Close()
	}()

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(file)

	// Parse SCV file into array of Payments
	paymentsArray, totalAmount, err := Handlers.ParseCSVFile(csvReader, configuration.Multiplier)
	if err != nil {
		log.Fatal(err)
	}

	if totalAmount != 6000000 {
		t.Errorf("Total payments amount in test_small_2.csv is %d, got: %f", 6000000, totalAmount)
	}

	if len(paymentsArray) != 3 {
		t.Errorf("We expect %d payments in test_small_2.csv, got: %d", 3, len(paymentsArray))
	}

	if paymentsArray[2].Amount.Quantity != 2500000 {
		t.Errorf("We expect %d amount at line %d in test_small_2.csv, got: %f", 2500000, 3, paymentsArray[2].Amount.Quantity)
	}

	if paymentsArray[2].Address != "addr_test1vpcsfa78jy5qwr40kzpr2s3ky7ga68al0qcsxgd9rkwnrecjvehy7" {
		t.Errorf("We expect %s address at line %d in test_small_2.csv, got: %s", "\"addr_test1vpcsfa78jy5qwr40kzpr2s3ky7ga68al0qcsxgd9rkwnrecjvehy7\"", 3, paymentsArray[2].Address)
	}
}
