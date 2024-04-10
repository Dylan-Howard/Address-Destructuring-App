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
	AddressUnit
	UnitValue string
	Delimiter string
}

func (list AddressList) Split() {
	values := strings.Split(list.UnitValue, list.Delimiter)
	fmt.Println(values)
}

type AddressRange struct {
	AddressUnit
	UnitDescriptor string
	StartValue 		 string
	EndValue 			 string
}

type AddressUnit struct {
	UnitDescriptor string
	UnitValue  		 string
}

type Address struct {
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

func (a Address) ToString() string {
	return a.StreetNumber + "," + a.StreetName + "," + a.Unit + "," + a.City
}
