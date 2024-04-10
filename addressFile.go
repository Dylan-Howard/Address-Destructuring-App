package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type AddressRow struct {
	StreetNumber	string
	StreetName		string
	Unit					string
	City					string
	Zip						string
	State					string
	Region				SchoolRegion
	InCounty			bool

}

type AddressFile struct {
	Rows []AddressRow
}

func (a AddressRow) ToString() string {
	return a.StreetNumber + "," + a.StreetName + "," + a.Unit + "," + a.City
}

func FetchAddressesFromCSV(filePath string, patterns Patterns, onlyWCPS bool) AddressFile {

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

	addressFile := new(AddressFile)

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

		addressFile.Rows = append(addressFile.Rows, tAddress)
	}

	return *addressFile
}

func generateUnitList(descriptor string, start int, end int, toString bool, prefix string) []string {
	var list = []string {}

	for i := start; i <= end; i++ {
		var val string

		if toString {
			val = string(rune(i))
		} else {
			val = strconv.Itoa(i)
		}

		unit := descriptor + " " + prefix +  val

		list = append(list, unit)
	}

	return list
}

func (a AddressRow) EnumerateRange(descriptor string, startValue string, endValue string) []string {
	// Determines if the range is character or integer based
	_, err := strconv.Atoi(startValue)

	var start, end int
	var isChar bool
	var prefix string

	// Sets the prefix if range uses character prefix before numbers
	if err != nil && len(startValue) > 1 {
		_, charErr := strconv.Atoi(string(startValue[1:]))
		if charErr == nil {
			prefix = string(startValue[0])
			startValue = startValue[1:]
		}
		if len(endValue) > 1 {
			endValue = endValue[1:]
		}
	}

	/* Sets parameters for generating the unit list */ 
	if err != nil {
		isChar = true

		start = int(startValue[0])
		end = int(endValue[0])
	} else {
		isChar = false

		start, _ = strconv.Atoi(startValue)
		end, _ = strconv.Atoi(endValue)
	}

	return generateUnitList(descriptor, start, end, isChar, prefix)
}

func (a AddressFile) ListAddresses(patterns Patterns) {
	list := []Address {}

	for i := 0; i < len(a.Rows); i++ {
		matchedPattern, isMatch := patterns.GetMatch(a.Rows[i].Unit)

		// if !isMatch {
		// 	fmt.Println("Error: no match found")
		// }

		if isMatch && matchedPattern.Type == "range" {
			// fmt.Println("Matched!")
			// fmt.Println(a.Rows[i].Unit)
			// fmt.Println(matchedPattern.Type)
			// fmt.Println(matchedPattern.Map)
			groups := matchedPattern.GetCaptureGroups(a.Rows[i].Unit)
			// fmt.Println(groups)

			descriptor := groups[matchedPattern.Map.Descriptor - 1]
			startValue := groups[matchedPattern.Map.StartValue - 1]
			endValue := groups[matchedPattern.Map.EndValue - 1]

			units := a.Rows[i].EnumerateRange(descriptor, startValue, endValue)

			for j := 0; j < len(units); j++ {
				var addr Address
				addr.StreetNumber = a.Rows[i].StreetNumber
				addr.StreetName = a.Rows[i].StreetName
				addr.City = a.Rows[i].City
				addr.Zip = a.Rows[i].Zip
				addr.State = a.Rows[i].State
				addr.Region = a.Rows[i].Region
				addr.InCounty = a.Rows[i].InCounty

				addr.Unit = units[j]

				list = append(list, addr)
			}
		}
	}

	for i := 0; i < len(list); i++ {
		println(list[i].ToString())
	}
}

func (a AddressFile) ExportAddressesToCSV(filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(a.Rows); i ++ {
		_, err = f.WriteString(a.Rows[i].ToString() + "\n")

		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}
}