package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func fetchAddressCollectionFromCSV(filePath string, onlyWCPS bool) AddressCollection {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	// Checks for the error
	if err != nil {
		fmt.Println("Error reading records")
	}

	addressFile := new(AddressCollection)

	for i := 0; i < len(records); i++ {
		var tAddress AddressRow
		tAddress.StreetNumber = strings.TrimSpace(records[i][3])
		tAddress.StreetName = strings.TrimSpace(records[i][4])
		tAddress.Unit = strings.TrimSpace(records[i][5])
		tAddress.City = strings.TrimSpace(records[i][6])
		tAddress.Zip = strings.TrimSpace(records[i][7])
		tAddress.State = strings.TrimSpace(records[i][12])

		tRegion := SchoolRegion{
			Elementary: strings.TrimSpace(records[i][16]),
			Middle: strings.TrimSpace(records[i][15]),
			High: strings.TrimSpace(records[i][14]),
		}

		tAddress.Region = tRegion

		if !onlyWCPS || tAddress.InCounty {
			addressFile.Rows = append(addressFile.Rows, tAddress)
		}
	}

	return *addressFile
}

func fetchSettings(filePath string) ConnectionSettings {
	/* Get Settings */
	settingsFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println("Couldn't open `settings.json`")
		fmt.Println(err)
	}
	defer settingsFile.Close()

	byteValue, _ := io.ReadAll(settingsFile)

	var settings ConnectionSettings
	json.Unmarshal(byteValue, &settings)

	return settings
}

func fetchQueries(filePath string) DatabaseQueries {
	/* Get Settings */
	queriesFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println("Couldn't open `queries.json`")
		fmt.Println(err)
	}
	defer queriesFile.Close()

	byteValue, _ := io.ReadAll(queriesFile)

	var queries DatabaseQueries
	json.Unmarshal(byteValue, &queries)

	return queries
}

func exportAddresses(filePath string, addresses []Address) {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(addresses); i ++ {
		_, err = f.WriteString(addresses[i].ToString() + "\n")

		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}
}

func main() {
	curDir, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	dataDirectory := filepath.Join(curDir, "data")

	/* Get patterns */
	patternsPath := filepath.Join(dataDirectory, "patterns.json")
	patterns := fetchPatternsFromJSON(patternsPath)

	/* Get addresses */
	addressesPath := filepath.Join(dataDirectory, "import", "wczp-addresses.csv")
	addressRecords := fetchAddressCollectionFromCSV(addressesPath, false)
	addresses := addressRecords.GetAddresses(patterns, "")

	/* Export Addresses */
	exportPath := filepath.Join(dataDirectory, "export", "addresses.csv")
	exportAddresses(exportPath, addresses)

	var campus InfiniteCampus

	settingsPath := filepath.Join(curDir, "settings", "campusSettings.json")
	campus.Settings = fetchSettings(settingsPath)

	queriesPath := filepath.Join(curDir, "settings", "campusQueries.json")
	campus.Queries = fetchQueries(queriesPath)

	campusAddresses := campus.GetAddresses()
}
