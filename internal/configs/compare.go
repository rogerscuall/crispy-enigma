package configs

import (
	"regexp"
	"sort"
	"strings"
)

// CompareNetworkConfigs compares two network configurations and returns true if they are functionally equivalent.
func CompareNetworkConfigs(running, designed string) bool {
	runningLines := processConfig(running)
	designedLines := processConfig(designed)

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
		return false
	}
	for i := range running {
		if running[i] != designed[i] {
			return false
		}
	}
	return true
}

