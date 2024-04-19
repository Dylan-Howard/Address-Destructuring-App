package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func fetchSettings() ConnectionSettings {
	/* Get Settings */
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	settingsPath := filepath.Join(curDir, "settings", "campusSettings.json")
	settingsFile, err := os.Open(settingsPath)

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

func fetchQueries() DatabaseQueries {
	/* Get Queries */
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	queriesPath := filepath.Join(curDir, "settings", "campusQueries.json")
	queriesFile, err := os.Open(queriesPath)

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

func getExportDirectory() string {
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	return filepath.Join(curDir, "data", "export")
}

func padStringLeft(str string, length int) string {
	for len(str) < length {
		str = "0" + str
	}
	return str
}

func getNewFileName(prefix string) string {
	exportDirectory := getExportDirectory()

	entries, err := os.ReadDir(exportDirectory)
  if err != nil {
    fmt.Println(err)
  }

	/* Find used filenames */
	usedFiles := make(map[string]bool)
  for _, e := range entries {
		usedFiles[e.Name()] = true
  }

	/* Generate unique filename */
	i := 1
	for {
		fileId := strconv.Itoa(i)
		tFilename := prefix + padStringLeft(fileId, 6) + ".csv"

		if !usedFiles[tFilename] {
			return tFilename
		}

		i++
	}
}

func ExportAddressesToCsv(prefix string, addresses []Address) (string, error) {

	/* Creates new file with unique name in the export directory */
	fileName := getNewFileName(prefix)
	exportFilePath := filepath.Join(getExportDirectory(), fileName)
	f, err := os.Create(exportFilePath)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("could not access the filepath")
	}

	/* Writes addresses to the new file */
	f.WriteString("Id,StreetNumber,StreetName,Unit,City,Zip,State,Region,InCounty\n")
	for i := 0; i < len(addresses); i ++ {
		_, err = f.WriteString(addresses[i].ToCsvString() + "\n")

		if err != nil {
			fmt.Println(err)
			f.Close()
			return "", errors.New("could write to the local directory")
		}
	}

	url := filepath.Join("data", "export", fileName)

	return url, nil
}

func getDifferences(addr1 []Address, addr2 []Address) ([]Address, []Address) {

	var toAdd []Address
	// var toRemove []Address

	for i := 0; i < len(addr1); i ++ {

		doesMatch := false
		j := 0

		for !doesMatch && j < len(addr2) {
			doesMatch = addr1[i].Match(addr2[j])
			j++
		}

		/* If not matched, queue this address to add */
		if j == len(addr2) {
			toAdd = append(toAdd, addr1[i])
		}

		/* Remove any matched addresses */
		if doesMatch {
			addr2[j] = addr2[len(addr2)-1]
    	addr2 = addr2[:len(addr2)-1]
		}
	}

	return toAdd, addr2
}

func executeCommit(validationId string) (int, int) {

	/* Fetch Commit Addresses */
	fmt.Println("Processing commit")
	validationPath := filepath.Join(getExportDirectory(), validationId + ".csv")
	commitAddresses := FetchAddressesFromLocalData(validationPath)
	
	/* Initialize Database Collection */
	var campus InfiniteCampus
	campus.Settings = fetchSettings()
	campus.Queries = fetchQueries()
	
	fmt.Println("Getting Campus Addresses")
	campusAddresses := campus.GetAddresses()
	
	fmt.Println("Processing Differences")
	addQueue, removeQueue := getDifferences(commitAddresses, campusAddresses)
	
	fmt.Println("Processing Add Queue")
	addCount, _ := campus.AddAddresses(addQueue)
	fmt.Println("Processing Remove Queue")
	removeCount, _ := campus.RemoveAddresses(removeQueue)

	return addCount, removeCount
}

func main() {
	fmt.Println("Your application is now running on: http://localhost:3000");

	Serve();
}
