package configs

import (
	"sort"
	"strings"
)

// CompareNetworkConfigs compares two network configurations and returns true if they are the same.
// This function ignores whitespace and order differences.
func CompareNetworkConfigs(running, designed string) bool {
	// Split the configurations into lines
	runningLines := strings.Split(running, "\n")
	designedLines := strings.Split(designed, "\n")

	// Trim whitespace and remove empty lines
	runningLines = cleanAndSort(runningLines)
	designedLines = cleanAndSort(designedLines)

	// Compare the sorted and cleaned lines
	return strings.Join(runningLines, "\n") == strings.Join(designedLines, "\n")
}

// cleanAndSort trims whitespace, removes empty lines, and sorts the resulting slice
func cleanAndSort(lines []string) []string {
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	sort.Strings(result)
	return result
}

