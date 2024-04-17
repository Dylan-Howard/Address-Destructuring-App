package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	fmt.Println("Your application is now running on: http://localhost:3000");

	Serve();
}
