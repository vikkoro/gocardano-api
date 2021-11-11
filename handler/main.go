package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"

	Wallets "github.com/vikkoro/gocardano-api/wallets"
)

// One endpoint for all requests
func (client ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get parameter from url request
	vars := mux.Vars(r)
	module := vars["module"]

	// Set our response header
	w.Header().Set("Content-Type", "application/json")

	// Handle each request using the module parameter:
	switch module {

	case "home":
		w.Header().Set("Content-Type", "text/html")

		fp := path.Join("templates", "home.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Get Network Information
	case "information":
		resp, err := http.Get(client.Configuration.InformationUrl)
		if err != nil {
			log.Fatalln(err)
		}

		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		sb := string(body)
		log.Printf(sb)

		json.NewEncoder(w).Encode(sb)

		// Get wallets list
	case "wallets":
		wallets := Wallets.GetList(client.Configuration.WalletsUrl)

		fmt.Println(wallets)
		json.NewEncoder(w).Encode(wallets)

		// Upload SCV file with addresses and amounts to be transferred
		/* SCV format we are expecting
		addr_test100000000000000000000000000000000000000000000000000001,1.956444
		addr_test100000000000000000000000000000000000000000000000000002,1.845180
		addr_test100000000000000000000000000000000000000000000000000003,2.395366
		*/
	case "upload":
		fmt.Println("File Upload Endpoint Hit")

		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(10 << 20)
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, handler, err := r.FormFile("csvFile")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		var pp []Wallets.PaymentData

		// Read csv values using csv.Reader
		csvReader := csv.NewReader(file)

		// We stockpile all transaction amounts here
		var aa float64

		for {
			// Read CSV file line by line
			rec, err := csvReader.Read()

			// Until the end of file
			if err == io.EOF {
				break
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal(err)
			}

			// We receive amount in CSV as float. Later we will need
			// to multiply it by 1M to use in the node API
			a, err := strconv.ParseFloat(rec[1], 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal(err)
			}

			// Use multiplier for compatibility with the node API
			// Add up transaction amounts
			aa += a * client.Configuration.Multiplier

			pp = append(pp, Wallets.PaymentData{
				Address: rec[0],
				Amount: Wallets.AmountData{
					Quantity: a * client.Configuration.Multiplier,
					Unit:     "lovelace",
				},
			})

			// do something with read line
			fmt.Printf("%+v\n", rec)
		}

		// Structure to check fees and send transactions
		transactions := &Wallets.TransactionsData{
			Passphrase: client.Configuration.Passphrase,
			Payments:   pp,
			TimeToLive: Wallets.AmountData{
				Quantity: 500,
				Unit:     "second",
			},
		}

		// Check our wallet current amount
		wallet := Wallets.GetWallet(client.Configuration.WalletsUrl, client.Configuration.WalletId)

		// Get estimated fees for the transactions
		estimated := Wallets.GetTransactionFee(client.Configuration.WalletsUrl, client.Configuration.WalletId, *transactions)

		// Check if wallet amount is enough
		if wallet.Balance.Available.Quantity < estimated.EstimatedMax.Quantity+aa {
			log.Fatal("Not enough funds")
		}

		// Send transactions
		response := Wallets.SendTransaction(client.Configuration.WalletsUrl, client.Configuration.WalletId, *transactions)

		// Create a JSON file to keep record of the transactions
		tempFile, err := ioutil.TempFile(client.Configuration.UploadDirectory, "transactions-*.json")
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()

		// Write transactions to our temporary file
		json.NewEncoder(tempFile).Encode(transactions)

		// Create a JSON file to keep record of the transactions
		respFile, err := ioutil.TempFile(client.Configuration.UploadDirectory, "response-*.json")
		if err != nil {
			fmt.Println(err)
		}
		defer respFile.Close()

		// Write transactions to our temporary file
		json.NewEncoder(respFile).Encode(response)

		// Return that we have successfully uploaded our file!
		fmt.Println(w, "Successfully Uploaded File\n")

		fmt.Println(wallet)

	}

}
