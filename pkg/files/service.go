package files

import (
	"encoding/json"
	"github.com/vikkoro/gocardano-api/pkg/config"
	"io/ioutil"
	"log"
	"os"
)

// Service interface used to list the strings
type Service interface {
	SaveJSONFile(string, interface{}) (string, error)
	ReadJSONFile(string, interface{}) (interface{}, error)
}

type service struct {
	cfg *config.Configuration
}

// NewService constructor of the default service.
func NewService(cfg *config.Configuration) *service {
	return &service{cfg}
}

// Parse SCV file into array of Payments
func (s *service) SaveJSONFile(dir string, content interface{}) (string, error) {

	// Create a JSON file to keep record of the transfers
	tempFile, err := ioutil.TempFile(dir, "transfers-*.json")
	if err != nil {
		log.Printf("FILES ERROR: %q", err.Error())
		return "", err
	}

	defer func() {
		_ = tempFile.Close()
	}()

	// Write transfers to our temporary file
	if err = json.NewEncoder(tempFile).Encode(content); err != nil {
		log.Printf("FILES ERROR: %q", err.Error())
		return "", err
	}

	return tempFile.Name(), nil
}

// Parse SCV file into array of Payments
func (s *service) ReadJSONFile(fileName string, structure interface{}) (interface{}, error) {

	file, err := os.Open(fileName)

	if err != nil {
		log.Printf("FILES ERROR: %q", err.Error())
		return nil, nil
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := json.NewDecoder(file)
	content := structure

	err = decoder.Decode(&content)
	if err != nil {
		log.Printf("FILES ERROR: %q", err.Error())
		return nil, nil
	}

	return &content, nil
}
