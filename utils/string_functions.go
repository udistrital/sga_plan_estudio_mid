// String functions.
// Custom functions that operate on strings.

package utils

import (
	"encoding/json"
	"strings"
)

// SplitTrimSpace use Split function to slices s into all substrings separated
// by sep and use TrimSpace to remove space and return a slice of the substrings.
func SplitTrimSpace(s, sep string) []string {
	substrings := strings.Split(s, sep)

	for i, elementString := range substrings {
		substrings[i] = strings.TrimSpace(elementString)
	}
	return substrings
}

func Map2String(inMap map[string]any) (string, error) {
	if outString, err := json.Marshal(&inMap); err == nil {
		return string(outString), nil
	} else {
		return "", err
	}
}
