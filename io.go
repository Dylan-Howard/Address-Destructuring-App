package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func fetchAddressesFromCSV(filePath string, patterns Patterns, onlyWCPS bool) ([]AddressRow) {

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

	var rows []AddressRow

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

		tType := patterns.ClassifyUnitValue(strings.TrimSpace(records[i][5]))

		fmt.Println(tType);

		if (onlyWCPS && tAddress.IsInDistrict()) || !onlyWCPS {
			rows = append(rows, tAddress)
		}
	}

	return rows
}

func ExportAddressesToCSV(filePath string, addresses []AddressRow) {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(addresses); i ++ {
		_, err = f.WriteString(addresses[i].ToString() + "\n")
		// _, err = f.WriteString(addresses[i].Unit + "\n")
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}
}