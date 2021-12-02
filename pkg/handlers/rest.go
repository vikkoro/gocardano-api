package handlers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vikkoro/gocardano-api/pkg/config"
	"github.com/vikkoro/gocardano-api/pkg/files"
	"github.com/vikkoro/gocardano-api/pkg/parser"
	"github.com/vikkoro/gocardano-api/pkg/wallet"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
)

func NewRestService(_c *config.Configuration, _w wallet.Service, _p parser.Service, _f files.Service) {
	rest := &restService{cfg: _c, w: _w, p: _p, f: _f}

	r := gin.Default()
	r.Static("/assets", "./assets")
	r.GET("/home", rest.getHome)
	r.GET("/api/health", rest.getHealth)
	r.GET("/api/v1/cardano/wallets/:id", rest.getWalletById)
	r.GET("/api/v1/cardano/wallets", rest.getWallets)
	r.POST("/api/v1/cardano/wallets/transfer", rest.transfer)

	_ = r.Run()
}

type restService struct {
	cfg *config.Configuration
	w   wallet.Service
	p   parser.Service
	f   files.Service
}

func (rs *restService) getHome(c *gin.Context) {

	err := Render(c.Writer, "home.html")
	if err != nil {
		SendError(c, err, http.StatusInternalServerError)
		return
	}
}

func (rs *restService) getHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "I'm Alive",
	})
}

func (rs *restService) getWallets(c *gin.Context) {
	sl, err := rs.w.GetWallets()
	if err != nil {
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{
		"status":  "OK",
		"message": sl,
	})
}

func (rs *restService) getWalletById(c *gin.Context) {
	id := c.Param("id")

	w, err := rs.w.GetWallet(id)
	if err != nil {
		SendError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(200, gin.H{
		"status":  "OK",
		"message": w,
	})
}

func (rs *restService) transfer(c *gin.Context) {

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		SendError(c, err, http.StatusBadRequest)
		return
	}
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := c.Request.FormFile("csvFile")
	if err != nil {
		SendError(c, err, http.StatusBadRequest)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		SendError(c, err, http.StatusBadRequest)
		return
	}

	payments, totalAmount, err := rs.p.ParsePayments(string(buf.Bytes()))
	if err != nil {
		SendError(c, err, http.StatusBadRequest)
		return
	}

	response, err := rs.w.Transfer(payments, totalAmount)
	if err != nil {
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{
		"status":  "OK",
		"message": response,
	})

	if _, err = rs.f.SaveJSONFile(rs.cfg.UploadDirectory, payments); err != nil {
		log.Printf("REST ERROR: %q", err.Error())
	}

}

// Render HTML page from template file
func Render(w http.ResponseWriter, templateFile string) error {
	w.Header().Set("Content-Type", "text/html")

	fp := path.Join("templates", templateFile)

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, ViewData{}); err != nil {
		return err
	}

	return nil
}

func SendError(c *gin.Context, err error, code uint64) {
	log.Printf("REST ERROR: %q", err.Error())

	switch code {

	case 400:
		c.JSON(400, gin.H{
			"status":  "BadRequest",
			"message": err.Error(),
		})

	case 500:
		c.JSON(500, gin.H{
			"status":  "InternalServerError",
			"message": err.Error(),
		})
	}

}