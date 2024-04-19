package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func fetchAddressCollectionFromData(rows [][]string, onlyWCPS bool) AddressCollection {
	addressFile := new(AddressCollection)

	for i := 0; i < len(rows); i++ {
		var tAddress AddressRow
		tAddress.Id, _ = strconv.Atoi(rows[i][0])
		tAddress.StreetNumber = strings.TrimSpace(rows[i][2])
		tAddress.StreetName = strings.TrimSpace(rows[i][3])
		tAddress.Unit = strings.TrimSpace(rows[i][5])
		tAddress.City = strings.TrimSpace(rows[i][6])
		tAddress.Zip = strings.TrimSpace(rows[i][8])
		tAddress.State = strings.TrimSpace(rows[i][12])
		tAddress.InCounty = strings.TrimSpace(rows[i][16]) != "Outside County Line"

		tRegion := SchoolRegion{
			Elementary: strings.TrimSpace(rows[i][15]),
			Middle: strings.TrimSpace(rows[i][14]),
			High: strings.TrimSpace(rows[i][13]),
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

func main() {
	fmt.Println("Your application is now running on: http://localhost:3000");

	Serve();
}
