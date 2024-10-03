package configs

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// CompareNetworkConfigs compares two network configurations and returns true if they are functionally equivalent.
func CompareNetworkConfigs(running, designed string) bool {
	runningLines := processConfig(running)
	designedLines := processConfig(designed)

	fmt.Println("Processed running config has", len(runningLines), "lines")
	fmt.Println("Processed designed config has", len(designedLines), "lines")
	return compareLines(runningLines, designedLines)
}

func processConfig(config string) []string {
	lines := strings.Split(config, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if shouldIncludeLine(line) {
			result = append(result, normalizeLine(line))
		}
	}
	sort.Strings(result)
	return result
}

func shouldIncludeLine(line string) bool {
	if line == "" || line == "!" {
		return false
	}
	excludePatterns := []string{
		`^username .* secret `,
		`^ *exec /usr/bin/TerminAttr `,
		`^ip route vrf MGMT`,
		`^ntp server vrf MGMT`,
	}
	for _, pattern := range excludePatterns {
		if matched, _ := regexp.MatchString(pattern, line); matched {
			return false
		}
	}
	return true
}

func normalizeLine(line string) string {
	// Remove trailing comments
	line = regexp.MustCompile(`\s*!.*$`).ReplaceAllString(line, "")

	// Normalize spaces
	line = regexp.MustCompile(`\s+`).ReplaceAllString(line, " ")

	return strings.TrimSpace(line)
}

func compareLines(running, designed []string) bool {
	if len(running) != len(designed) {
		fmt.Printf("Length mismatch: running has %d lines, designed has %d lines\n", len(running), len(designed))
		return false
	}
	for i := range running {
		if running[i] != designed[i] {
			fmt.Printf("Difference at line %d:\nRunning:  %s\nDesigned: %s\n", i+1, running[i], designed[i])
			return false
		}
	}
	return true
}

// // CompareNetworkConfigs compares two network configurations and returns true if they are functionally equivalent.
// func CompareNetworkConfigs(running, designed string) bool {
// 	runningLines := processConfig(running)
// 	designedLines := processConfig(designed)

// 	fmt.Printf("Processed running config has %d lines\n", len(runningLines))
// 	fmt.Printf("Processed designed config has %d lines\n", len(designedLines))

// 	return compareLines(runningLines, designedLines)
// }
