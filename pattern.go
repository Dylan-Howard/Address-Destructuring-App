package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
)

// Pattern struct
type Pattern struct {
	Expression string `json:"expression"`
	Type   	 string `json:"type"`
}

// Patterns struct
type Patterns struct {
	Patterns []Pattern `json:"patterns"`
}


func fetchPatternsFromJSON(filePath string) (Patterns) {

	/* Get Patterns */
	patternFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
			fmt.Println(err)
	}
	fmt.Println("Successfully opened `patterns.json`")
	// defer the closing of our jsonFile so that we can parse it later on
	defer patternFile.Close()

	// json to Byte array
	byteValue, _ := io.ReadAll(patternFile)

	var patterns Patterns

	json.Unmarshal(byteValue, &patterns)

	return patterns
}

func countOperations(patterns Patterns, addresses []AddressRow) (map[int]int) {
	counts := make(map[int]int)

	for i := 1; i < len(addresses); i++ {
		for j := len(patterns.Patterns) - 1; j >= 0; j-- {
			var exp, _ = regexp.Compile(patterns.Patterns[j].Expression)

			matched := exp.MatchString(addresses[i].Unit)

			if (matched) {
				counts[j] = counts[j] + 1
				break
			}
			
		}
	}

	return counts
}

func SortRows(patterns Patterns, addresses []AddressRow) [][]AddressRow {
	var sorted [][]AddressRow

	for j := 0; j <= len(patterns.Patterns); j ++ {
		var group []AddressRow
		sorted = append(sorted, group)
	}

	for i := 1; i < len(addresses); i++ {
		for j := len(patterns.Patterns) - 1; j >= 0; j-- {
			var exp, _ = regexp.Compile(patterns.Patterns[j].Expression)

			matched := exp.MatchString(addresses[i].Unit)

			if matched {
				// fmt.Println(exp.FindStringSubmatch(addresses[i].Unit));
				fmt.Println(patterns.Patterns[j].Type)
				sorted[j] = append(sorted[j], addresses[i])
				break
			}

			if j > 0 {
				continue
			}

			sorted[len(sorted)-1] = append(sorted[len(sorted)-1], addresses[i])
		}
	}

	return sorted
}

func (p Patterns) ClassifyUnitValue(value string) string {

	for j := len(p.Patterns) - 1; j >= 0; j-- {
		var exp, _ = regexp.Compile(p.Patterns[j].Expression)

		matched := exp.MatchString(value)

		if matched {
			// fmt.Println(exp.FindStringSubmatch(addresses[i].Unit));
			return p.Patterns[j].Type
		}
	}

	return "unknown"
}

func (p Pattern) MatchParts(address AddressRow) ([]string) {
	var expression, regErr = regexp.Compile(p.Expression)
	if regErr != nil {
		return nil
	}

	matches := expression.FindStringSubmatch(address.Unit)

	if (len(matches) > 1) {
		return matches[1:]
	}

	return matches
}