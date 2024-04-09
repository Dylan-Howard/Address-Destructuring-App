package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	curDir, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	dataDirectory := filepath.Join(curDir, "data")

	/* Get patterns */
	patternsPath := filepath.Join(dataDirectory, "patterns.json")
	fmt.Println(patternsPath)
	patterns := fetchPatternsFromJSON(patternsPath)
	fmt.Println(len(patterns.Patterns))

	/* Get addresses */
	addressesPath := filepath.Join(dataDirectory, "import", "wczp-addresses.csv")
	addressRecords := fetchAddressesFromCSV(addressesPath, patterns, false)
	fmt.Println(len(addressRecords));
		
	/* Loop to iterate through and print each of the string slice */
	counts := countOperations(patterns, addressRecords)

	for i := 0; i < len(patterns.Patterns); i++ {
		fmt.Printf("%d addresses matched pattern %d\n", counts[i], i + 1);
	}

	/* Loop to iterate through and print each of the string slice */
	groups := SortRows(patterns, addressRecords)
	fmt.Printf("%d addresses were unmatched\n", len(groups[len(groups)-1]));

	for i := 0; i < len(groups); i ++ {
		groupFileName := "group" + strconv.Itoa(i) + ".csv"
		testPath := filepath.Join(dataDirectory, "export", groupFileName)
		ExportAddressesToCSV(testPath, groups[i])
	}
	
}
