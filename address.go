package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
 * Address Struct - Used as the output form of address destructuring
 */
type SchoolRegion struct {
	Elementary	string
	Middle			string
	High				string
}

type Address struct {
	Id						int
	StreetNumber	string
	StreetName		string
	Unit					string
	City					string
	Zip						string
	State					string
	Region				SchoolRegion
	InCounty			bool
}

func (a Address) IsInDistrict() bool {
	return a.Region.Elementary != "City Schools"
}

func (a Address) ToCsvString() string {
	return "\"" + (strconv.Itoa(a.Id) + "\",\"" + a.StreetNumber + "\",\"" + a.StreetName + "\",\"" + a.Unit + "\",\"" +
		a.City + "\",\"" + a.Zip + "\",\"" + a.State + "\",\"[" + a.Region.Elementary + ";" +
		a.Region.Middle + ";" + a.Region.High + "]\",\"" + strconv.FormatBool(a.InCounty)) + "\""
}

func (a Address) ToString() string {
	return a.StreetNumber + "," + a.StreetName + "," + a.Unit + "," + a.City
}

func (a1 Address) Match(a2 Address) bool {
	return a1.StreetNumber == a2.StreetNumber && a1.StreetName == a2.StreetName && a1.Unit == a2.Unit
}

func (a1 Address) Compare(a2 Address) string {
	comparisonResults := []string{}
	comparisonResults = append(comparisonResults, "The following fields do not match: ")
	if a1.Id != a2.Id {
		comparisonResults = append(comparisonResults, "Id")
	}
	if a1.StreetNumber != a2.StreetNumber {
		comparisonResults = append(comparisonResults, "StreetNumber")
	}
	if a1.StreetName != a2.StreetName {
		comparisonResults = append(comparisonResults, "StreetName")
	}
	if a1.Unit != a2.Unit {
		comparisonResults = append(comparisonResults, "Unit")
	}
	if a1.City != a2.City {
		comparisonResults = append(comparisonResults, "City")
	}
	if a1.State != a2.State {
		comparisonResults = append(comparisonResults, "State")
	}
	if a1.Zip != a2.Zip {
		comparisonResults = append(comparisonResults, "Zip")
	}
	if a1.Region != a2.Region {
		comparisonResults = append(comparisonResults, "Region")
	}
	if a1.InCounty != a2.InCounty {
		comparisonResults = append(comparisonResults, "InCounty")
	}

	return strings.Join(comparisonResults, ", ")
}

/*
 * Address Row - A single address data row. Many compose an Address Collection
 */
type AddressRow struct {
	Id						int
	StreetNumber	string
	StreetName		string
	Unit					string
	City					string
	Zip						string
	State					string
	Region				SchoolRegion
	InCounty			bool
}

func (a AddressRow) ToString() string {
	return a.StreetNumber + "," + a.StreetName + "," + a.Unit + "," + a.City
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

func (a AddressRow) EnumerateRange(pattern Pattern) []Address {
	list := []Address {}

	groups := pattern.GetCaptureGroups(a.Unit)

	descriptor := groups[pattern.Map.Descriptor - 1]
	startValue := groups[pattern.Map.StartValue - 1]
	endValue := groups[pattern.Map.EndValue - 1]

	/* Determines if the range is character or integer based */
	_, err := strconv.Atoi(startValue)

	var start, end int
	var isChar bool
	var prefix string
	
	/* Sets the prefix if range uses character prefix before numbers */
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
	
	units := generateUnitList(descriptor, start, end, isChar, prefix)

	/* Initializes Addresses for each unit */
	for j := 0; j < len(units); j++ {
		var addr Address
		addr.Id = a.Id
		addr.StreetNumber = a.StreetNumber
		addr.StreetName = a.StreetName
		addr.City = a.City
		addr.Zip = a.Zip
		addr.State = a.State
		addr.Region = a.Region
		addr.InCounty = a.InCounty

		addr.Unit = units[j]

		list = append(list, addr)
	}

	return list
}

func (a AddressRow) Split(pattern Pattern) AddressCollection {
	var collection AddressCollection

	groups := pattern.GetCaptureGroups(a.Unit)
	delimiter := groups[pattern.Map.Delimiter - 1]
	splitItems := strings.Split(groups[0], delimiter)

	for i := 0; i < len(splitItems); i++ {
		var row AddressRow
		row.Id = a.Id
		row.StreetNumber = a.StreetNumber
		row.StreetName = a.StreetName
		row.Unit = strings.TrimSpace(splitItems[i])
		row.City = a.City
		row.Zip = a.Zip
		row.State = a.State
		row.Region = a.Region
		row.InCounty = a.InCounty

		collection.Rows = append(collection.Rows, row)
	}

	return collection
}

func (a AddressRow) BuildUnit(pattern Pattern) []Address {
	var list []Address

	groups := pattern.GetCaptureGroups(a.Unit)

	var descriptor, value string
	if pattern.Map.Descriptor != 0 {
		descriptor = groups[pattern.Map.Descriptor - 1] + " "
	}
	if pattern.Map.UnitValue != 0 {
		value = groups[pattern.Map.UnitValue - 1]
	}

	var addr Address
	addr.Id = a.Id
	addr.StreetNumber = a.StreetNumber
	addr.StreetName = a.StreetName
	addr.City = a.City
	addr.Zip = a.Zip
	addr.State = a.State
	addr.Region = a.Region
	addr.InCounty = a.InCounty

	addr.Unit = descriptor + value

	list = append(list, addr)
	
	return list
}

/*
 * Address Collection - A collection of many address rows.
 */
type AddressCollection struct {
	Rows []AddressRow
}

func (a AddressCollection) GetAddresses(patterns Patterns, typeFilter string) []Address {
	list := []Address {}

	for i := 0; i < len(a.Rows); i++ {
		matchedPattern, isMatch := patterns.GetMatch(a.Rows[i].Unit)

		var items []Address

		if isMatch && matchedPattern.Type == "range" {
			items = a.Rows[i].EnumerateRange(matchedPattern)
		}
		if isMatch && matchedPattern.Type == "list" {
			collection := a.Rows[i].Split(matchedPattern)
			items = collection.GetAddresses(patterns, "")
		}
		if isMatch && matchedPattern.Type == "unit" {
			items = a.Rows[i].BuildUnit(matchedPattern)
		}

		if typeFilter == "" || matchedPattern.Type == typeFilter {
			list = append(list, items...)
		}
	}

	return list
}

func BuildAddressCollectionFromData(rows [][]string, onlyWCPS bool) AddressCollection {
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

func FetchAddressesFromLocalData(filePath string) []Address {
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

	addresses := []Address{}

	for i := 1; i < len(records); i++ {
		// fmt.Println(records[i])

		var tAddress Address
		tAddress.Id, _ = strconv.Atoi(records[i][0])
		tAddress.StreetNumber = records[i][1]
		tAddress.StreetName = records[i][2]
		tAddress.Unit = records[i][3]
		tAddress.City = records[i][4]
		tAddress.Zip = records[i][5]
		tAddress.State = records[i][6]
		tAddress.InCounty = records[i][8] == "true"

		regionParts := strings.Split(records[i][7], ";")
		tRegion := SchoolRegion{
			Middle: regionParts[1],
		}
		if 0 < len(regionParts[0]) {
			tRegion.Elementary = regionParts[0][1:]
		}
		if (0 < len(regionParts[2])) {
			tRegion.High = regionParts[2][0:len(regionParts[2])-1]
		}

		tAddress.Region = tRegion

		addresses = append(addresses, tAddress)
	}

	return addresses
}
