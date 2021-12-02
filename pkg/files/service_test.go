package files

import (
	"github.com/vikkoro/gocardano-api/pkg/config"
	"testing"
)

var cfg *config.Configuration

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

type Content struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func TestService_SaveJSONFile(t *testing.T) {
	// Given
	s := NewService(cfg)

	fileContent := &Content{Title: "Test Title", Message: "Test Message"}

	// When
	fileName, err := s.SaveJSONFile(s.cfg.UploadDirectory, fileContent)

	// Then
	if err != nil {
		t.Errorf("Not error expected, but got: %q", err.Error())
	}

	// When reading content from the file
	_, err = s.ReadJSONFile(fileName, Content{})

	// Then
	if err != nil {
		t.Errorf("Not error expected, but got: %q", err.Error())
	}

	// Cast read structure to Content
	//cc := c.(*Content)
	//
	//if err == nil &&  cc.Title != "Test Title"{
	//	t.Errorf("Did not get expected result. Wanted: %q, got: %q", "Test Title", cc.Title)
	//}
}
