package main

import (
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
	return (strconv.Itoa(a.Id) + "," + a.StreetNumber + "," + a.StreetName + "," + a.Unit + "," +
		a.City + "," + a.Zip + "," + a.State + ",[" + a.Region.Elementary + ";" +
		a.Region.Middle + ";" + a.Region.High + "]," + strconv.FormatBool(a.InCounty))
}

func (a Address) ToString() string {
	return a.StreetNumber + "," + a.StreetName + "," + a.Unit + "," + a.City
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
