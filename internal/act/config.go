package act

import (
	"bufio"
	"fmt"
	"strings"
)

// BlockMatcher defines the interface for matching blocks
type BlockMatcher interface {
	Match(line string) bool
	IsEnd(line string) bool
}

// BlockUpdater defines the interface for updating blocks
type BlockUpdater interface {
	UpdateBlock(block []string) []string
}

type GenericMatcher struct {
	// Keyword is the string that identifies the beginning of a block
	Keyword string
}

// NewGenericMatcher creates a new GenericMatcher with the given keyword
func NewGenericMatcher(keyword string) GenericMatcher {
	return GenericMatcher{Keyword: keyword}
}

func (m GenericMatcher) Match(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), m.Keyword)
}

func (m GenericMatcher) IsEnd(line string) bool {
	return strings.TrimSpace(line) == "!"
}

// SingleLineUpdater modifies or removes a matched single line
type SingleLineUpdater struct {
	NewLine string // If NewLine is empty, the line will be removed
}

// UpdateBlock replaces the line with the new line
// If NewLine is empty, the line will be removed
func (u SingleLineUpdater) UpdateBlock(block []string) []string {
	if u.NewLine == "" {
		// If NewLine is empty, return an empty slice to remove the line
		return []string{}
	}
	// Replace the line with the new line
	return []string{u.NewLine, "!"}
}

// NestedMatcher matches lines inside a block (e.g., "ip address, mtu" inside an interface block)
type NestedMatcher struct {
	ParentMatcher BlockMatcher
	Keyword       string
	inParentBlock bool
}

func (m *NestedMatcher) Match(line string) bool {
	if m.ParentMatcher.Match(line) {
		m.inParentBlock = true
		return false
	}

	if m.inParentBlock && strings.Contains(strings.TrimSpace(line), m.Keyword) {
		return true
	}

	if m.inParentBlock && m.ParentMatcher.IsEnd(line) {
		m.inParentBlock = false
	}

	return false
}

func (m NestedMatcher) IsEnd(line string) bool {
	return m.ParentMatcher.IsEnd(line)
}

func NewNestedMatcher(parent BlockMatcher, keyword string) *NestedMatcher {
	return &NestedMatcher{ParentMatcher: parent, Keyword: keyword}
}

// GenericProcessor processes the configuration based on matchers and updaters
// It will apply the first matching matcher and updater
type GenericProcessor struct {
	Matchers []BlockMatcher
	Updaters []BlockUpdater
}

/*
ProcessConfig processes the configuration line by line
If a line matches a matcher, the updater will be applied once it reaches the end of the block
*/
func (p *GenericProcessor) ProcessConfig(config string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(config))
	var currentBlock []string
	inBlock := false
	currentMatcher := BlockMatcher(nil)
	currentUpdater := BlockUpdater(nil)

	for scanner.Scan() {
		line := scanner.Text()
		if inBlock {
			//3. If the line is in a block, it is added to the current block but not to the result
			currentBlock = append(currentBlock, line)
			if currentMatcher.IsEnd(line) {
				//4. If the line is the end of the block, the updater is applied to the block
				updatedBlock := currentUpdater.UpdateBlock(currentBlock)
				//5. The updated block is added to the result
				for _, updatedLine := range updatedBlock {
					result.WriteString(updatedLine + "\n")
				}
				//6. The current block is reset
				inBlock = false
				currentBlock = nil
				currentMatcher = nil
				currentUpdater = nil
			}
		} else {
			//2. Each line is checked against all the matchers to see if it is the beginning of a block
			// TODO: Only the first matching matcher is applied. Add support for multiple matchers
			for i, matcher := range p.Matchers {
				if matcher.Match(line) {
					inBlock = true
					currentMatcher = matcher
					currentUpdater = p.Updaters[i]
					currentBlock = append(currentBlock, line)
					break
				}
			}
			//1. Most of the time this line is used to append the line to the result for those lines that are not in the block
			if !inBlock {
				result.WriteString(line + "\n")
			}
		}
	}

	return result.String()
}

// func clean() {

// }

// BlockUpdaterFunc is a helper to create BlockUpdater from a function
type BlockUpdaterFunc func([]string) []string

func (f BlockUpdaterFunc) UpdateBlock(block []string) []string {
	return f(block)
}

var usernames []string = []string{
	"username arista privilege 15 role network-admin secret sha512 $6$ZGX/X07MoiWP9hvX$3UaAtOAiBGc54DYHdQt5dsr5P2HLydxjrda51Zw69tSsF4tahXPVj26PwOiZUy/xFRZL3CAkT7.lsOPqWfIbU0",
	"username cvpadmin secret sha512 $6$vO.NVE0FD54oFstS$IyyRB7.D/30GF2jg89Ep9HcqTlmR0jx.gCoipN8cHbhK..U7qZ2dzDjcEg56cmb5L78jNPq6y3yviJkr48vCc0",
	"username ec2-user shell /bin/bash nopassword",
	"username ec2-user ssh-key ssh-rsa ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7zrvZRnQAErF4WL1z5HRDzbG2OGxCliQ2GvgmVOsO1fhlnvnzQkd2or/aMO58YK/ytBCraV7pm+zpdgzckkdOeJGCWBPps5V8MltIGaJvNM2Wp0FakjyLv7k+4p/rGPzbJ5gRd9sF1d+5dJ4WhuWd8d+182+snLZF+TapQD/jYRFnoacUDFLShZVvvVkMJpr9ENm5i+Rt1e1MqWTNHct47oKeAmQ9bkkX/EWS4dnDMRHxmb2IC+u2vUlfGeAGOI9LMMaDKZ67KpT9LyTk4H0hVQzKM16OO/IcnFafGG4EQSkxPRkDGESngELid1EhVKbK8pMkY7VJg6t5QW7Cf5bT20cl8aHzM/r9buFQ7k9Qfgq6Kr33kXOuxONiKQB2396oXNiuidynHXtqmK+hYwlvt+BqOJRvfaNdiCM//Aa6ZPSQ/mdJhb9brwQ5wY/ITxaQoHBDRrfIHHWw4s0l/tr2RBfEXg/5bSLf4DbUH7A7fvMWZ2Wl783q4mHFXqtmx/8= root@buildkitsandbox",
	"username service shell /bin/bash secret sha512 $6$YIifMegwrRBSVaGk$31svzVjkGQwhAX4QPwtBUpN.WLKmFc7qDFNdCLguC7zaJ3Mn2oATcoIUKQAUm32YdQNTxYZc8091YaOI4yxa71",
}

var dameonChange string = "    exec /usr/bin/TerminAttr -cvcompression=gzip -smashexcludes=ale,flexCounter,hardware,kni,pulse,strata -ingestexclude=/Sysdb/cell/1/agent,/Sysdb/cell/2/agent -cvaddr=%s -cvauth=token,/tmp/token -cvvrf=%s -taillogs"
var daemonNew []string = []string{
	"daemon TerminAttr",
	"    no shutdown",
	"!"}

func CleanConfig(config, cvpHost, vrf string) string {
	// Define matchers and updaters
	interfaceMatcherMTU := NewGenericMatcher("interface")
	daemonMatcher := NewGenericMatcher("daemon")
	//Match the block queue-monitor streaming "qms"
	qmsMatcher := NewGenericMatcher("queue-monitor streaming")
	mtuMatcher := NewNestedMatcher(interfaceMatcherMTU, "mtu")
	singleLineMatcher := NewGenericMatcher("username")
	singleLineMatcherMonitor := NewGenericMatcher("monitor")
	singleLineMatcherAAA := NewGenericMatcher("aaa")
	singleLineSNMP := NewGenericMatcher("snmp-server")
	singleLineQueueMonitorLenght := NewGenericMatcher("queue-monitor length")

	// Define updaters

	mtuUpdater := BlockUpdaterFunc(func(block []string) []string {
		for i, line := range block {
			if strings.Contains(strings.TrimSpace(line), "mtu") {
				block[i] = "   mtu 1500"
			}
		}
		return block
	})

	daemonUpdate := fmt.Sprintf(dameonChange, cvpHost, vrf)
	position := 1
	daemonNew := append(daemonNew[:position], append([]string{daemonUpdate}, daemonNew[position:]...)...)
	singleLineUpdater := SingleLineUpdater{NewLine: strings.Join(usernames, "\n")}
	singleLineUpdaterRemover := SingleLineUpdater{NewLine: ""}
	daemonUpdater := SingleLineUpdater{NewLine: strings.Join(daemonNew, "\n")}
	qmsUpdated := SingleLineUpdater{NewLine: "!"}

	// Create the processor with matchers and updaters
	processor := GenericProcessor{
		Matchers: []BlockMatcher{mtuMatcher, singleLineMatcher, singleLineMatcherMonitor, daemonMatcher, singleLineMatcherAAA, singleLineSNMP, singleLineQueueMonitorLenght, qmsMatcher},
		Updaters: []BlockUpdater{mtuUpdater, singleLineUpdater, singleLineUpdaterRemover, daemonUpdater, singleLineUpdaterRemover, singleLineUpdaterRemover, singleLineUpdaterRemover, qmsUpdated},
	}

	// Process the configuration
	return processor.ProcessConfig(config)
}

type ActCVP struct {
	Host string
	VRF  string
	Port string
}

type ActCVPConfigs []ActCVP
