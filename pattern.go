package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
)

// Pattern struct
type PatternMap struct {
	UnitValue  int	`json:"unitValue"`
	Descriptor int	`json:"descriptor"`
	StartValue int	`json:"startValue"`
	EndValue   int	`json:"endValue"`
	Delimiter  int	`json:"delimiter"`
}

// Pattern struct
type Pattern struct {
	Expression string 		`json:"expression"`
	Type   	   string 		`json:"type"`
	Map   	   PatternMap `json:"map"`
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
		fmt.Println("Couldn't open `patterns.json`")
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer patternFile.Close()

	// json to Byte array
	byteValue, _ := io.ReadAll(patternFile)

	var patterns Patterns

	json.Unmarshal(byteValue, &patterns)

	return patterns
}

func CountOperations(patterns Patterns, addresses []AddressRow) (map[int]int) {
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

func (p Patterns) GetMatch(value string) (Pattern, bool) {

	pattern := Pattern{}

	for i := len(p.Patterns) - 1; i >= 0; i-- {
		var exp, _ = regexp.Compile(p.Patterns[i].Expression)

		matched := exp.MatchString(value)

		if matched {
			return p.Patterns[i], true
		}
	}

	return pattern, false
}

func (p Pattern) GetCaptureGroups(value string) []string {
	var exp, _ = regexp.Compile(p.Expression)

	return exp.FindStringSubmatch(value)
}