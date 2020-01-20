package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/falcinspire/hackpackpdf/internal/app"
)

var lastRequest time.Time

func FileUpload(w http.ResponseWriter, req *http.Request) {

	now := time.Now()
	elapsed := now.Sub(lastRequest)
	if elapsed.Seconds() < 10 {
		http.Error(w, "Server was just busy... taking a short breath", http.StatusTooManyRequests)
		return
	}

	dir, err := ioutil.TempDir("", "project")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	acceptZip(w, req)
	processZip()
	respondPdf(w)
	err = os.Chdir(wd)
	if err != nil {
		panic(err)
	}

	lastRequest = time.Now()
}

func acceptZip(w http.ResponseWriter, req *http.Request) error {
	// Source: https://dzone.com/articles/go-archive-support-for-microservice
	// TODO check license
	if req.Method == "POST" {
		err := req.ParseMultipartForm(32 * 1024)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		file, _, err := req.FormFile("uploadfile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		defer file.Close()
		size, err := file.Seek(0, 2)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		file.Seek(0, 0)
		unzipWeb(file, size, ".")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		return nil
	} else {
		http.Error(w, "Must be POST request", http.StatusBadRequest)
		return fmt.Errorf("Must be POST request")
	}
}

func processZip() {
	err := app.CompileHackpack("hackpack.toml", "hackpack.pdf")
	if err != nil {
		panic(err)
	}
}

func respondPdf(w http.ResponseWriter) {
	//Check if file exists and open
	hackpackpdf, err := os.Open("hackpack.pdf")
	defer hackpackpdf.Close() //Close after function return
	if err != nil {
		panic(err)
	}

	stats, _ := hackpackpdf.Stat()                 //Get info from file
	pdfsize := strconv.FormatInt(stats.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename=hackpack.pdf")
	w.Header().Set("Content-Type", "pdf")
	w.Header().Set("Content-Length", pdfsize)

	io.Copy(w, hackpackpdf) //'Copy' the file to the client
}

func WebService() {
	lastRequest = time.Now()
	http.HandleFunc("/upload", FileUpload)
	http.ListenAndServe(":80", nil)
}
