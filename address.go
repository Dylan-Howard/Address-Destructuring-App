package main

import (
	"fmt"
	"strings"
)

type SchoolRegion struct {
	Elementary	string
	Middle			string
	High				string
}

type AddressList struct {
	UnitValue string
	Delimiter string
}

func (list AddressList) Split() {
	values := strings.Split(list.UnitValue, list.Delimiter)
	fmt.Println(values)
}

type AddressRange struct {
	UnitValue string
	Includer string
}

type AddressUnit struct {
	Unit string
}

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

func (row AddressRow) IsInDistrict() bool {
	return row.Region.Elementary != "City Schools"
}

func (row AddressRow) ToString() string {
	return row.StreetNumber + "," + row.StreetName + "," + row.Unit + "," + row.City
}
