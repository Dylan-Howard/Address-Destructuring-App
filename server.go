package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type AddressUploadFile struct {
	Headers []string `json:"headers"`
	Rows 		[][]string `json:"rows"`
}

type ValidationResponse struct {
  Data []Address `json:"data"`
  Status int `json:"status"`
}

type SubmitResponse struct {
	Count  int `json:"count"`
  Status int `json:"status"`
}

func handleAddressValidation(rw http.ResponseWriter, req *http.Request) {

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

	fmt.Println("Headers:", len(addr.Headers))
  fmt.Println("Rows:", len(addr.Rows))

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
	fmt.Println(addressRecords.Rows[0])
	addresses := addressRecords.GetAddresses(patterns, "")

	fmt.Println("Addresses:", len(addresses))

	responseData := ValidationResponse{
		Data: addresses,
		Status: 200,
	}
	
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		return
	}
	
	rw.Write(jsonData)
}

func handleAddressSubmit(rw http.ResponseWriter, req *http.Request) {

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
	http.Handle("/", fs)
	http.HandleFunc("/api/addresses/validate", handleAddressValidation)
	http.HandleFunc("/api/addresses/submit", handleAddressSubmit)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
