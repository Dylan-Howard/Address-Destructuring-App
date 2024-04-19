package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AddressUploadFile struct {
	Headers []string 	 `json:"headers"`
	Rows 		[][]string `json:"rows"`
}

type ValidationResponse struct {
  Data 	 []Address `json:"data"`
	Url		 string		 `json:"url"`
  Status int 			 `json:"status"`
}

type SubmitResponse struct {
	Count  int `json:"count"`
  Status int `json:"status"`
}

func handleValidation(rw http.ResponseWriter, req *http.Request) {

	header := rw.Header()
	header.Add("Content-Type", "application/json")
	header.Add("Access-Control-Allow-Origin", "http://localhost:3000")
	header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	if req.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodPost {
    return
  }
  
  /* Decode the request body into the user struct */
	addr := &AddressUploadFile{}

  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(addr)
  if err != nil {
    return
  }

	/* Get patterns */
	curDir, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	dataDirectory := filepath.Join(curDir, "./data")
	patternsPath := filepath.Join(dataDirectory, "patterns.json")
	patterns := fetchPatternsFromJSON(patternsPath)
  
  /* Process address data */
	addressRecords := fetchAddressCollectionFromData(addr.Rows, false)
	addresses := addressRecords.GetAddresses(patterns, "")

	/* Create saved file */
	validationFile, err := ExportAddressesToCsv("validation", addresses)
	if err != nil {
		log.Println(err)
	}

	responseData := ValidationResponse{
		Data: addresses,
		Status: 200,
		Url: validationFile,
	}
	
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		return
	}
	
	rw.Write(jsonData)
}

func handleCommit(rw http.ResponseWriter, req *http.Request) {

	fmt.Println("Handling commit")

	header := rw.Header()
	header.Add("Content-Type", "application/json")
	header.Add("Access-Control-Allow-Origin", "http://localhost:3000")
	header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	if req.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodPost {
    return
  }

	responseData := SubmitResponse{
		Count: 10,
		Status: 200,
	}
	
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		return
	}

	rw.Write(jsonData)
}

func Serve() {
	fs := http.FileServer(http.Dir("analuo/build"))
	directory := flag.String("d", ".", "the directory of static file to host")
	http.Handle("/", fs)
	http.Handle("/data/export/", http.StripPrefix(strings.TrimRight("/data/export/", "/"), http.FileServer(http.Dir(*directory))))
	http.HandleFunc("/api/validation", handleValidation)
	http.HandleFunc("/api/commit", handleCommit)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
