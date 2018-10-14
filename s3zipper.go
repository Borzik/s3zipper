package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"net/http"
)

type configuration struct {
	Port   int
	BindTo string
}

var config = configuration{}

func main() {
	configFile, _ := os.Open("conf.json")
	decoder := json.NewDecoder(configFile)
	err := decoder.Decode(&config)
	if err != nil {
		panic("Error reading conf")
	}

	fmt.Println("Running on port", config.Port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(config.BindTo+":"+strconv.Itoa(config.Port), nil)
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	strs, ok := r.URL.Query()["token"]
	if !ok || len(strs) < 1 {
		http.Error(w, "File Zipper. Pass ?token= to use.", 403)
		return
	}
	str := strs[0]

	eventID, _ := r.URL.Query()["event_id"]
	host, _ := r.URL.Query()["host"]

	resp, err := httpClient.Get(fmt.Sprint("https://", host[0], "/api/events/", eventID[0], "/photos/s3zipper?token=", str))
	if err != nil {
		http.Error(w, err.Error(), 403)
		log.Printf("%s\t%s\t%s", r.Method, r.RequestURI, err.Error())
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "error", 403)
		log.Printf("%s\t%s", r.Method, r.RequestURI)
		return
	}

	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		http.Error(w, err2.Error(), 403)
		log.Printf("%s\t%s\t%s", r.Method, r.RequestURI, err2.Error())
		return
	}
	var arr []string
	_ = json.Unmarshal(bodyBytes, &arr)

	w.Header().Add("Content-Disposition", fmt.Sprint("attachment; filename=", resp.Header.Get("Zip-File-Name")))
	w.Header().Add("Content-Type", "application/zip")

	zipWriter := zip.NewWriter(w)

	for _, file := range arr {
		var rdr io.ReadCloser
		var err error

		var res *http.Response
		res, err = http.Get(file)

		if err != nil {
			log.Printf("Error loading \"%s\" - %s", file, err.Error())
			continue
		}
		rdr = res.Body

		name := path.Base(res.Request.URL.Path)
		h := &zip.FileHeader{Name: name, Method: zip.Deflate, Flags: 0x800}
		f, err := zipWriter.CreateHeader(h)

		if err != nil {
			log.Printf("Error adding \"%s\" - %s", file, err.Error())
			continue
		}

		io.Copy(f, rdr)
		rdr.Close()
	}

	zipWriter.Close()

	log.Printf("%s\t%s\t%s", r.Method, r.RequestURI, time.Since(start))
}
