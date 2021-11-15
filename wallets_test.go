package main

import (
	"fmt"
	Handlers "github.com/vikkoro/gocardano-api/handler"
	Wallets "github.com/vikkoro/gocardano-api/wallets"
	"log"
	"os"
	"testing"
)

func TestGenerateAddress(t *testing.T) {
	configuration := GetConfig("conf.json")

	file, _ := os.Open("./data/test_small_2.csv")

	defer file.Close()

	// Parse SCV file into array of Payments
	pp, totalAmount, err := Handlers.ParseCSVFile(file, configuration.Multiplier)
	if err != nil {
		log.Fatal(err)
	}

	// Structure to check fees and send transactions
	transactions := &Wallets.TransactionsData{
		Passphrase: configuration.Passphrase,
		Payments:   pp,
		TimeToLive: Wallets.AmountData{
			Quantity: 500,
			Unit:     "second",
		},
	}

	// Check our wallet current amount
	wallet := Wallets.GetWallet(configuration.WalletsUrl, configuration.WalletId)
	balanceBefore := wallet.Balance.Available.Quantity

	// Get estimated fees for the transactions
	estimated := Wallets.GetTransactionFee(configuration.WalletsUrl, configuration.WalletId, *transactions)

	// Check if wallet amount is enough
	if wallet.Balance.Available.Quantity < estimated.EstimatedMax.Quantity+totalAmount {
		log.Fatal("Not enough funds")
	}

	// Send transactions
	Wallets.SendTransaction(configuration.WalletsUrl, configuration.WalletId, *transactions)

	wallet = Wallets.GetWallet(configuration.WalletsUrl, configuration.WalletId)
	balanceAfter := wallet.Balance.Available.Quantity

	fmt.Printf("Amount %f \n", totalAmount)
	fmt.Printf("Before %f \n", balanceBefore)
	fmt.Printf("After %f \n", balanceAfter)
	fmt.Printf("Fee min %f \n", estimated.EstimatedMin.Quantity)
	fmt.Printf("Fee max %f \n", estimated.EstimatedMax.Quantity)

	//t.Errorf("Sum was incorrect, got: %d, want: %d.", 33, 10)

}
