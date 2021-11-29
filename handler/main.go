package handlers

import (
	"encoding/csv"
	"encoding/json"
	. "errors"
	"fmt"
	"github.com/gorilla/mux"
	Wallets "github.com/vikkoro/gocardano-api/wallets"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
)

// One endpoint for all requests
func (client ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get parameter from url request
	vars := mux.Vars(r)
	module := vars["module"]

	// Set our response header
	w.Header().Set("Content-Type", "application/json")

	err := ProcessRequest(module, client.Configuration, w, r)

	if err != nil {
		fmt.Println("handler:", err)

		if err = json.NewEncoder(w).Encode(&Error{Code: 500, Message: err.Error()}); err != nil {
			fmt.Println("handler:", err)
		}
	}
}

// Process requests to different modules
func ProcessRequest(module string, config Configuration, w http.ResponseWriter, r *http.Request) error {
	// Handle each request using the module parameter:
	switch module {

	case "home":

		return Render(w, "home.html", ViewData{})

	// Get Network Information
	case "information":

		response, err := SendGETRequest(config.InformationUrl)
		if err != nil {
			return err
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			return err
		}

	// Get wallets list
	case "wallets":

		wallets, err := Wallets.GetList(config.WalletsUrl)
		if err != nil {
			return err
		}

		fmt.Println(wallets)

		if err = json.NewEncoder(w).Encode(wallets); err != nil {
			return err
		}

	// Upload SCV file with addresses and amounts to be transferred
	/* SCV format we are expecting
	addr_test100000000000000000000000000000000000000000000000000001,1.956444
	addr_test100000000000000000000000000000000000000000000000000002,1.845180
	addr_test100000000000000000000000000000000000000000000000000003,2.395366
	*/
	case "upload":
		fmt.Println("File Upload Endpoint Hit")

		response, err := UploadCSVFile(config, r)
		if err != nil {
			return err
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			return err
		}
	}

	return nil
}

// Upload SCV file
func UploadCSVFile(config Configuration, r *http.Request) (*Wallets.TransferFundsResponseData, error) {

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, err
	}
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("csvFile")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(file)

	// Parse SCV file into array of Payments
	pp, totalAmount, err := ParseCSVFile(csvReader, config.Multiplier)
	if err != nil {
		return nil, err
	}

	// Structure to check fees and send transfers
	payments := &Wallets.BulkPaymentsData{
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
		return nil, err
	}

	// Get estimated fees for the transfers
	estimated, err := Wallets.GetTransferFee(config.WalletsUrl, config.WalletId, *payments)
	if err != nil {
		return nil, err
	}

	// Check if wallet amount is enough
	if wallet.Balance.Available.Quantity < estimated.EstimatedMax.Quantity+totalAmount {
		return nil, New("not enough funds")
	}

	// Send a transfer
	response, err := Wallets.TransferFunds(config.WalletsUrl, config.WalletId, *payments)
	if err != nil {
		return nil, err
	}

	// Create a JSON file to keep record of the transfers
	tempFile, err := ioutil.TempFile(config.UploadDirectory, "transfers-*.json")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tempFile.Close()
	}()

	// Write transfers to our temporary file
	if err = json.NewEncoder(tempFile).Encode(payments); err != nil {
		return nil, err
	}

	return response, nil
}

// Parse SCV file into array of Payments
func ParseCSVFile(csvReader *csv.Reader, multiplier uint64) ([]Wallets.PaymentData, float64, error) {

	var paymentsArray []Wallets.PaymentData

	// We stockpile all transfers amounts here
	var totalAmountLovelace float64

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
			return nil, 0, New(fmt.Sprintf("transfer amount is lower than 1 ADA at line %d", lineCounter))
		}

		// Use multiplier for compatibility with the node API
		// Add up transfers amounts
		totalAmountLovelace += amountAda * float64(multiplier)

		paymentsArray = append(paymentsArray, Wallets.PaymentData{
			Address: rec[0],
			Amount: Wallets.AmountData{
				Quantity: amountAda * float64(multiplier),
				Unit:     "lovelace",
			},
		})

		// do something with read line
		fmt.Printf("%+v\n", rec)

		lineCounter++
	}

	return paymentsArray, totalAmountLovelace, nil
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
func Render(w http.ResponseWriter, templateFile string, data ViewData) error {
	w.Header().Set("Content-Type", "text/html")

	fp := path.Join("templates", templateFile)

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, data); err != nil {
		return err
	}

	return nil
}
