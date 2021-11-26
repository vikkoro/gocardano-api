package handlers

import (
	"encoding/csv"
	"encoding/json"
	. "errors"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
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
		Render(w, "home.html", ViewData{})

	// Get Network Information
	case "information":
		response, err := SendGETRequest(client.Configuration.InformationUrl)

		if err != nil {
			fmt.Println(err)

			_ = json.NewEncoder(w).Encode(&Error{Code: 500, Message: err.Error()})
			return
		}

		_ = json.NewEncoder(w).Encode(response)

	// Get wallets list
	case "wallets":
		wallets, err := Wallets.GetList(client.Configuration.WalletsUrl)

		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)

			_ = json.NewEncoder(w).Encode(&Error{Code: 500, Message: err.Error()})
			return
		}

		fmt.Println(wallets)
		_ = json.NewEncoder(w).Encode(wallets)

	// Upload SCV file with addresses and amounts to be transferred
	/* SCV format we are expecting
	addr_test100000000000000000000000000000000000000000000000000001,1.956444
	addr_test100000000000000000000000000000000000000000000000000002,1.845180
	addr_test100000000000000000000000000000000000000000000000000003,2.395366
	*/
	case "upload":
		fmt.Println("File Upload Endpoint Hit")

		response, err := UploadCSVFile(client.Configuration, r)

		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)

			_ = json.NewEncoder(w).Encode(&Error{Code: 500, Message: err.Error()})
			return
		}

		_ = json.NewEncoder(w).Encode(response)
	}

}

// Upload SCV file
func UploadCSVFile(config Configuration, r *http.Request) (Wallets.SendTransactionsResponseData, error) {

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	_ = r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("csvFile")
	if err != nil {
		return Wallets.SendTransactionsResponseData{}, err
	}

	defer func() {
		_ = file.Close()
	}()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Parse SCV file into array of Payments
	pp, totalAmount, err := ParseCSVFile(file, config.Multiplier)
	if err != nil {
		return Wallets.SendTransactionsResponseData{}, err
	}

	// Structure to check fees and send transactions
	transactions := &Wallets.TransactionsData{
		Passphrase: config.Passphrase,
		Payments:   pp,
		TimeToLive: Wallets.AmountData{
			Quantity: 500,
			Unit:     "second",
		},
	}

	// Check our wallet current amount
	wallet, err := Wallets.GetWallet(config.WalletsUrl, config.WalletId)
	if err != nil {
		return Wallets.SendTransactionsResponseData{}, err
	}

	// Get estimated fees for the transaction
	estimated, err := Wallets.GetTransactionFee(config.WalletsUrl, config.WalletId, *transactions)
	if err != nil {
		return Wallets.SendTransactionsResponseData{}, err
	}

	// Check if wallet amount is enough
	if wallet.Balance.Available.Quantity < estimated.EstimatedMax.Quantity+totalAmount {
		return Wallets.SendTransactionsResponseData{}, New("not enough funds")
	}

	// Send transaction
	response, err := Wallets.SendTransaction(config.WalletsUrl, config.WalletId, *transactions)

	// Create a JSON file to keep record of the transactions
	tempFile, err := ioutil.TempFile(config.UploadDirectory, "transactions-*.json")
	if err != nil {
		fmt.Println(err)
	}

	defer func() {
		_ = tempFile.Close()
	}()

	// Write transactions to our temporary file
	_ = json.NewEncoder(tempFile).Encode(transactions)

	return response, nil
}

// Parse SCV file into array of Payments
func ParseCSVFile(file multipart.File, multiplier uint64) ([]Wallets.PaymentData, float64, error) {
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
			return nil, 0, err
		}

		// We receive amount in CSV as float. Later we will need
		// to multiply it by 1M to use in the node API
		a, err := strconv.ParseFloat(rec[1], 64)
		if err != nil {
			return nil, 0, err
		}

		// Use multiplier for compatibility with the node API
		// Add up transaction amounts
		aa += a * float64(multiplier)

		pp = append(pp, Wallets.PaymentData{
			Address: rec[0],
			Amount: Wallets.AmountData{
				Quantity: a * float64(multiplier),
				Unit:     "lovelace",
			},
		})

		// do something with read line
		fmt.Printf("%+v\n", rec)
	}

	return pp, aa, nil
}

func SendGETRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	sb := string(body)
	log.Printf(sb)

	return sb, nil
}

// Render HTML page from template file
func Render(w http.ResponseWriter, templateFile string, data ViewData) {
	w.Header().Set("Content-Type", "text/html")

	fp := path.Join("templates", templateFile)
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
