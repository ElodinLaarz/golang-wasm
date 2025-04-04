package main

import (
	"encoding/json"
	"sort"
	"strings"
	"syscall/js"
)

type ItemData struct {
	Name      string              `json:"name"`
	Variables map[string][]string `json:"variables"`
}

type FuzzyResult struct {
	Results []string `json:"results"`
}

var database = map[string]ItemData{
	"Apple": {
		Name: "Apple",
		Variables: map[string][]string{
			"Color": {"Red", "Green"},
			"Taste": {"Sweet", "Tart"},
		},
	},
	"Banana": {
		Name: "Banana",
		Variables: map[string][]string{
			"Color":   {"Yellow"},
			"Texture": {"Soft", "Creamy"},
		},
	},
	"Orange": {
		Name: "Orange",
		Variables: map[string][]string{
			"Color": {"Orange"},
			"Taste": {"Citrus", "Tangy"},
		},
	},
	"Mango": {
		Name: "Mango",
		Variables: map[string][]string{
			"Color":   {"Yellow", "Red"},
			"Texture": {"Smooth"},
		},
	},
	"Pineapple": {
		Name: "Pineapple",
		Variables: map[string][]string{
			"Color": {"Brown", "Yellow"},
			"Taste": {"Tangy", "Sweet"},
		},
	},
	"Grapes": {
		Name: "Grapes",
		Variables: map[string][]string{
			"Color":   {"Green", "Red", "Purple"},
			"Texture": {"Juicy"},
		},
	},
	"Strawberry": {
		Name: "Strawberry",
		Variables: map[string][]string{
			"Color": {"Red"},
			"Taste": {"Sweet", "Sour"},
		},
	},
	"Blueberry": {
		Name: "Blueberry",
		Variables: map[string][]string{
			"Color": {"Blue"},
			"Taste": {"Sweet"},
		},
	},
}

func fuzzyFind(pattern string) []string {
	pattern = strings.ToLower(pattern)
	var candidates []string
	for item := range database {
		candidates = append(candidates, item)
	}
	returnScores := map[string]int{}
	for _, candidate := range candidates {
		candidate = strings.ToLower(candidate)
		patternCharsMatched := 0
		fuzzyScore := 0 // higher is better, more skipped letters after first match is worse.
		for candidateIndex := 0; candidateIndex < len(candidate) && patternCharsMatched < len(pattern); candidateIndex++ {
			if candidate[candidateIndex] == pattern[patternCharsMatched] {
				patternCharsMatched++
				continue
			}
			if patternCharsMatched > 0 {
				fuzzyScore--
			}
		}
		if patternCharsMatched == len(pattern) {
			returnScores[candidate] = fuzzyScore
		}
	}
	sortedCandidates := make([]string, 0, len(returnScores))
	for candidate := range returnScores {
		sortedCandidates = append(sortedCandidates, candidate)
	}
	sort.Slice(sortedCandidates, func(i, j int) bool {
		if returnScores[sortedCandidates[i]] < returnScores[sortedCandidates[j]] {
			return true
		}
		return sortedCandidates[i] < sortedCandidates[j]
	})
	return sortedCandidates
}

// Why in the world is this the interface? ANY?! Helppppp!
func fuzzyFindJS(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "No pattern provided"
	}
	pattern := args[0].String()
	results := fuzzyFind(pattern)
	if len(results) == 0 {
		return "No matches found :("
	}
	fuzzyResult := &FuzzyResult{
		Results: results,
	}
	b, err := json.Marshal(fuzzyResult)
	if err != nil {
		return "Error encoding JSON"
	}
	return string(b)
}

func lookup(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "No item provided"
	}
	item := args[0].String()
	data, ok := database[item]
	if !ok {
		return "No data found for " + item
	}
	b, err := json.Marshal(data)
	if err != nil {
		return "Error encoding JSON"
	}
	return string(b)
}

func main() {
	js.Global().Set("lookupItem", js.FuncOf(lookup))
	js.Global().Set("fuzzyFind", js.FuncOf(fuzzyFindJS))
	select {}
}
