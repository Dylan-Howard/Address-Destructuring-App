
type SchoolRegion = {
  "Elementary": string,
  "Middle": string,
  "High": string,
};

export type Address = {
  "Id": number,
	"StreetNumber": string,
	"StreetName": string,
	"Unit": string,
	"City": string,
	"Zip": string,
	"State": string,
	"Region": SchoolRegion,
	"InCounty": boolean,
};