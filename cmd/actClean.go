/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// actCleanCmd represents the actClean command
var actCleanCmd = &cobra.Command{
	Use:   "actClean",
	Short: "Cleans a production AVD designed configuration to be used with ACT",
	Long: `Arista Cloud Test (ACT) is runs in virtualized devices, that do not support all the physical devices features.
This command will clean a production AVD designed configuration to be used with ACT.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("actClean called")
		clean()
	},
}

func init() {
	rootCmd.AddCommand(actCleanCmd)
	//TODO: we need to ask for the CVP IP address,
	//TODO: The IP of CVP is the one defined in the ACT topology, the port seems to be static in 9910, is the internal address not the public
   //TODO: we need to ask for the CVP VRF, the default is MGMT
}

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

func clean() {
	// Example configuration string

	// Define matchers and updaters
	interfaceMatcherMTU := NewGenericMatcher("interface")
	daemonMatcher := NewGenericMatcher("daemon")
	mtuMatcher := NewNestedMatcher(interfaceMatcherMTU, "mtu")
	singleLineMatcher := NewGenericMatcher("username")
	singleLineMatcherMonitor := NewGenericMatcher("monitor")
	singleLineMatcherAAA := NewGenericMatcher("aaa")

	// Define updaters

	mtuUpdater := BlockUpdaterFunc(func(block []string) []string {
		for i, line := range block {
			if strings.Contains(strings.TrimSpace(line), "mtu") {
				block[i] = "   mtu 1500"
			}
		}
		return block
	})

	daemonUpdate := fmt.Sprintf(dameonChange, "10.255.33.114:9910", "MGMT")
	position := 1
	daemonNew := append(daemonNew[:position], append([]string{daemonUpdate}, daemonNew[position:]...)...)
	singleLineUpdater := SingleLineUpdater{NewLine: strings.Join(usernames, "\n")}
	singleLineUpdaterRemover := SingleLineUpdater{NewLine: ""}
	daemoonUpdater := SingleLineUpdater{NewLine: strings.Join(daemonNew, "\n")}

	// Create the processor with matchers and updaters
	processor := GenericProcessor{
		Matchers: []BlockMatcher{mtuMatcher, singleLineMatcher, singleLineMatcherMonitor, daemonMatcher, singleLineMatcherAAA},
		Updaters: []BlockUpdater{mtuUpdater, singleLineUpdater, singleLineUpdaterRemover, daemoonUpdater, singleLineUpdaterRemover},
	}

	// Process the configuration
	updatedConfig := processor.ProcessConfig(config)
	fmt.Println(updatedConfig)
}

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

var config = `
!RANCID-CONTENT-TYPE: arista
!
daemon TerminAttr
   exec /usr/bin/TerminAttr -cvaddr=10.157.18.5:9910 -cvauth=certs,/persist/secure/ssl/terminattr/primary/certs/client.crt,/persist/secure/ssl/terminattr/primary/keys/client.key -cvvrf=MGMT -disableaaa -smashexcludes=ale,flexCounter,hardware,kni,pulse,strata -taillogs
   no shutdown
!
vlan internal order ascending range 1006 1199
!
transceiver qsfp default-mode 4x10G
!
service routing protocols model multi-agent
!
logging synchronous level critical
logging vrf MGMT host 10.167.12.233 514
logging vrf MGMT host 159.165.12.25 514
logging vrf MGMT host cribl-atl1.delta.com 514
!
hostname ATL-ADM-BL201
ip name-server vrf MGMT 159.165.56.240
dns domain delta.com
!
ntp server vrf MGMT ntp.delta.com prefer
!
sflow sample 16384
sflow vrf default destination 127.0.0.1
sflow vrf default source-interface Loopback0
sflow run
!
snmp-server contact nmc.dt@delta.com 404-714-6004
snmp-server location ATL,USA,GA
snmp-server vrf MGMT local-interface Management1
snmp-server group F1yDL03Rv3 v3 auth
snmp-server user admin admin v3 auth SHA password priv AES password
snmp-server host 159.165.13.22 vrf MGMT version 3 auth admin
snmp-server enable traps
snmp-server vrf MGMT
!
hardware speed-group 2 serdes 10g
hardware speed-group 4 serdes 10g
hardware speed-group 5 serdes 10g
hardware speed-group 7 serdes 10g
hardware speed-group 8 serdes 10g
!
spanning-tree mode mstp
spanning-tree mst 0 priority 4096
!
service unsupported-transceiver Delta 636f2557
!
tacacs-server host 159.165.34.251 vrf MGMT key 7 04690C25222C621C1F4E2A1514
!
aaa group server tacacs+ TACACS
   server 159.165.34.251 vrf MGMT
!
aaa authentication login default group TACACS local
aaa authentication enable default group TACACS local
aaa authorization serial-console
aaa authorization exec default group TACACS local
aaa authorization commands 1 default group TACACS local
aaa authorization commands 15 default group TACACS local
!
no enable password
no aaa root
!
username admin privilege 15 role network-admin secret sha512 $6$j7pR03FftRKAmvcI$ydCVF.2dSX1covKL.n3NwB0kVxiKixWGvtH5le2928aGbDGKOtLHbSVIxxdPZpAn8tmg2zauEVqDx86EDegvy1
username ngdc privilege 15 role network-admin secret sha512 $6$CuT89D1Uf1lIyUWB$BkqdlKWYcZ9k3wDE20DDr6u2kMtWE4U/.NTYO2.bmectBPJTjTB0502NAv0Bq4jT4k9spkp2Tlkvnz.qwBPBu/
!
vlan 21
   name ATL-ADM-FW-OUT
!
vlan 22
   name ATL-ADM-FW-IN
!
vlan 25
   name ATL-ADM-FW-SI-OUT
!
vlan 26
   name ATL-ADM-SI-FW-IN
!
vlan 1101
   name Mgmt-ESXI-V1101
!
vlan 1102
   name Vmotion-ESXI-V1102
!
vlan 1103
   name backup-V1103
!
vlan 1105
   name Vmotion-SNS-ESXI-V1105
!
vlan 1106
   name backup-SNS-V1106
!
vlan 1107
   name VM-MGMT-V1107
!
vlan 1108
   name dmf-tunnel
!
vlan 1109
   name example
!
vlan 1301
   name SNS_L2_GW_FW_10.157.96.1/24
!
vlan 1302
   name SNS_L2_GW_FW_10.157.97.1/24
!
vlan 1303
   name SNS_L2_GW_FW_10.157.98.1/24
!
vlan 1701
   name Mgmt-ESXI-V1701
!
vlan 1702
   name Vmotion-ESXI-V1702
!
vlan 1703
   name Backup-V1703
!
vlan 1704
   name VM-MGMT-V1704
!
vlan 1705
   name subnet-F-5
!
vlan 1749
   name SAP-TEST-NODE
!
vlan 1801
   name SNS_SI_L2_GW_FW_10.157.192.1/24
!
vlan 1802
   name SNS_SI_L2_GW_FW_10.157.193.1/24
!
vlan 1803
   name SNS_SI_L2_GW_FW_10.157.194.1/24
!
vlan 2001
   name Prod_DCC_App2
!
vlan 2002
   name Prod_DCC_DB2
!
vlan 2003
   name Prod_DCP_APP2
!
vlan 2004
   name Prod_DCP_DB1
!
vlan 2005
   name Prod_APP1
!
vlan 2006
   name Prod_NETAPP
!
vlan 2007
   name Prod_DCP_APP1
!
vlan 2008
   name Prod_DCC_DB1
!
vlan 2009
   name Prod_DB1
!
vlan 2010
   name Prod_DCC_APP1
!
vlan 2101
   name DEV_DCC_APP2
!
vlan 2102
   name DEV_DCC_DB2
!
vlan 2103
   name DEV_DCC_APP1
!
vlan 2104
   name DEV_DB1
!
vlan 2105
   name DEV_APP1
!
vlan 2106
   name Dev_DCC_DB1
!
vlan 2107
   name DEV_DCP_DB1
!
vlan 2108
   name DEV_DB1
!
vlan 2109
   name SI_DCC_APP1
!
vlan 2110
   name SI_DCC_DB2
!
vlan 2111
   name SI_Exadata
!
vlan 2112
   name SI_DCC_DB1
!
vlan 2113
   name SI_DB1
!
vlan 2114
   name SI_DCP_APP1
!
vlan 2115
   name SI_APP1
!
vlan 2116
   name SI_DCP_DB1
!
vlan 2117
   name SI_APP1
!
vlan 2118
   name SI_NETAPP
!
vlan 2217
   name SAP-NIM
!
vrf instance CORE
!
vrf instance DL_Prod
!
vrf instance DL_SI
!
vrf instance MGMT
!
interface Port-Channel251
   description FWATLADMINTP02A_AGG
   no shutdown
   switchport
   switchport trunk allowed vlan 2-2200
   switchport trunk native vlan 4092
   switchport mode trunk
   evpn ethernet-segment
      identifier 0000:0202:f111:d63d:15e1
      route-target import f1:11:d6:3d:15:e1
   lacp system-id f111.d63d.15e1
!
interface Port-Channel261
   description FWATLADMINTP02B_AGG
   no shutdown
   switchport
   switchport trunk allowed vlan 2-2200
   switchport trunk native vlan 4092
   switchport mode trunk
   evpn ethernet-segment
      identifier 0000:0202:dea2:2bcc:73b5
      route-target import de:a2:2b:cc:73:b5
   lacp system-id dea2.2bcc.73b5
!
interface Port-Channel271
   description FWATLADMINTS01A_AGG
   no shutdown
   switchport
   switchport trunk allowed vlan 2-2200
   switchport trunk native vlan 4092
   switchport mode trunk
   evpn ethernet-segment
      identifier 0000:0202:6fca:4678:917f
      route-target import 6f:ca:46:78:91:7f
   lacp system-id 6fca.4678.917f
!
interface Port-Channel281
   description FWATLADMINTS02B_AGG
   no shutdown
   switchport
   switchport trunk allowed vlan 2-2200
   switchport trunk native vlan 4092
   switchport mode trunk
   evpn ethernet-segment
      identifier 0000:0202:17f9:0836:0241
      route-target import 17:f9:08:36:02:41
   lacp system-id 17f9.0836.0241
!
interface Port-Channel291
   description ATL-ADM-PROD-LB-A_AGG
   no shutdown
   switchport
   switchport trunk allowed vlan 2-2200
   switchport trunk native vlan 4092
   switchport mode trunk
   evpn ethernet-segment
      identifier 0000:0202:7313:83e1:adad
      route-target import 73:13:83:e1:ad:ad
   lacp system-id 7313.83e1.adad
!
interface Port-Channel301
   description ATL-ADM-PROD-LB-B_AGG
   no shutdown
   switchport
   switchport trunk allowed vlan 2-2200
   switchport trunk native vlan 4092
   switchport mode trunk
   evpn ethernet-segment
      identifier 0000:0202:a566:2108:4ba5
      route-target import a5:66:21:08:4b:a5
   lacp system-id a566.2108.4ba5
!
interface Ethernet1/1
   description P2P_LINK_TO_ATL-ADM-SP101_Ethernet1/1
   no shutdown
   mtu 9214
   no switchport
   ip address 10.157.0.61/31
!
interface Ethernet2/1
   description P2P_LINK_TO_ATL-ADM-SP102_Ethernet1/1
   no shutdown
   mtu 9214
   no switchport
   ip address 10.157.0.63/31
!
interface Ethernet3/1
   description P2P_LINK_TO_ATL-ADM-SP103_Ethernet1/1
   no shutdown
   mtu 9214
   no switchport
   ip address 10.157.0.65/31
!
interface Ethernet4/1
   description SPAN-DST-ATL-ADM-PBN-LF201-eth1
   no shutdown
   switchport
!
interface Ethernet7/1
   description ATLADM5N-DC01:ETH8/5
   no shutdown
   mtu 9000
   no switchport
   vrf CORE
   ip address 10.157.13.0/31
!
interface Ethernet8/1
   description ATLADM5N-DC02:ETH8/5
   no shutdown
   mtu 9000
   no switchport
   vrf CORE
   ip address 10.157.13.2/31
!
interface Ethernet13/1
   speed 10g
   switchport
!
interface Ethernet15/1
   speed 10g
   switchport
!
interface Ethernet17/1
   speed 10g
   switchport
!
interface Ethernet19/1
   speed 10g
   switchport
!
interface Ethernet25/1
   description FWATLADMINTP02A_Ethernet41
   no shutdown
   channel-group 251 mode active
!
interface Ethernet26/1
   description FWATLADMINTP02B_Ethernet41
   no shutdown
   channel-group 261 mode active
!
interface Ethernet27/1
   description FWATLADMINTS01A_Ethernet41
   no shutdown
   channel-group 271 mode active
!
interface Ethernet28/1
   description FWATLADMINTS02B_Ethernet41
   no shutdown
   channel-group 281 mode active
!
interface Ethernet29/1
   description ATL-ADM-PROD-LB-A_Ethernet1
   no shutdown
   channel-group 291 mode active
!
interface Ethernet30/1
   description ATL-ADM-PROD-LB-B_Ethernet1
   no shutdown
   channel-group 301 mode active
!
interface Loopback0
   description EVPN_Overlay_Peering
   no shutdown
   ip address 10.157.4.11/32
!
interface Loopback1
   description VTEP_VXLAN_Tunnel_Source
   no shutdown
   ip address 10.157.6.11/32
!
interface Loopback100
   description vrf CORE
   no shutdown
   vrf CORE
   ip address 10.157.12.1/32
!
interface Loopback101
   description vrf DL_Prod
   no shutdown
   vrf DL_Prod
   ip address 10.157.12.3/32
!
interface Loopback201
   description vrf DL_SI
   no shutdown
   vrf DL_SI
   ip address 10.157.12.13/32
!
interface Management1
   description oob_management
   no shutdown
   vrf MGMT
   ip address 10.157.18.21/24
!
interface Vlan21
   description ATL-ADM-FW-OUT
   no shutdown
   vrf CORE
   ip address 10.157.14.2/29
!
interface Vlan22
   description ATL-ADM-FW-IN
   no shutdown
   vrf DL_Prod
   ip address 10.157.14.10/29
!
interface Vlan25
   description ATL-ADM-FW-SI-OUT
   no shutdown
   vrf CORE
   ip address 10.157.14.34/29
!
interface Vlan26
   description ATL-ADM-SI-FW-IN
   no shutdown
   vrf DL_SI
   ip address 10.157.14.42/29
!
interface Vlan1101
   description Mgmt-ESXI-V1101
   no shutdown
   vrf DL_Prod
   ip address virtual 10.157.48.1/24
!
interface Vlan1102
   description Vmotion-ESXI-V1102
   no shutdown
   vrf DL_Prod
   ip address virtual 10.157.49.1/24
!
interface Vlan1103
   description backup-V1103
   no shutdown
   vrf DL_Prod
   ip address virtual 10.157.50.1/24
!
interface Vlan1105
   description Vmotion-SNS-ESXI-V1105
   no shutdown
   vrf DL_Prod
   ip address virtual 10.157.52.1/24
!
interface Vlan1106
   description backup-SNS-V1106
   no shutdown
   vrf DL_Prod
   ip address virtual 10.157.53.1/24
!
interface Vlan1107
   description VM-MGMT-V1107
   no shutdown
   vrf DL_Prod
   ip address virtual 10.157.54.1/24
!
interface Vlan1108
   description dmf-tunnel
   no shutdown
   vrf DL_Prod
   ip address virtual 10.157.55.1/24
!
interface Vlan1109
   description example
   no shutdown
   vrf DL_Prod
   ip address virtual 10.157.56.1/24
!
interface Vlan1701
   description Mgmt-ESXI-V1701
   no shutdown
   vrf DL_SI
   ip address virtual 10.157.160.1/24
!
interface Vlan1702
   description Vmotion-ESXI-V1702
   no shutdown
   vrf DL_SI
   ip address virtual 10.157.161.1/24
!
interface Vlan1703
   description Backup-V1703
   no shutdown
   vrf DL_SI
   ip address virtual 10.157.162.1/24
!
interface Vlan1704
   description VM-MGMT-V1704
   no shutdown
   vrf DL_SI
   ip address virtual 10.157.163.1/24
!
interface Vlan1705
   description subnet-F-5
   no shutdown
   vrf DL_SI
   ip address virtual 10.157.164.1/24
!
interface Vlan1749
   description SAP-TEST-NODE
   no shutdown
   vrf DL_SI
   ip address virtual 10.157.140.1/27
!
interface Vlan2001
   description Prod_DCC_App2
   shutdown
   vrf DL_Prod
   ip address virtual 10.151.52.1/23
!
interface Vlan2002
   description Prod_DCC_DB2
   shutdown
   vrf DL_Prod
   ip address virtual 10.151.61.1/24
!
interface Vlan2003
   description Prod_DCP_APP2
   shutdown
   vrf DL_Prod
   ip address virtual 10.151.71.1/24
!
interface Vlan2004
   description Prod_DCP_DB1
   shutdown
   vrf DL_Prod
   ip address virtual 10.151.80.1/24
!
interface Vlan2005
   description Prod_APP1
   shutdown
   vrf DL_Prod
   ip address virtual 10.151.40.1/24
!
interface Vlan2006
   description Prod_NETAPP
   shutdown
   vrf DL_Prod
   ip address virtual 10.151.127.1/24
!
interface Vlan2007
   description Prod_DCP_APP1
   shutdown
   vrf DL_Prod
   ip address virtual 10.151.70.1/24
!
interface Vlan2008
   description Prod_DCC_DB1
   shutdown
   vrf DL_Prod
   ip address virtual 10.151.60.1/24
!
interface Vlan2009
   description Prod_DB1
   shutdown
   vrf DL_Prod
   ip address virtual 10.151.41.1/24
!
interface Vlan2010
   description Prod_DCC_APP1
   shutdown
   vrf DL_Prod
   ip address virtual 10.151.50.1/24
!
interface Vlan2101
   description DEV_DCC_APP2
   shutdown
   vrf DL_SI
   ip address virtual 10.151.174.1/23
!
interface Vlan2102
   description DEV_DCC_DB2
   shutdown
   vrf DL_SI
   ip address virtual 10.151.173.1/24
!
interface Vlan2103
   description DEV_DCC_APP1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.170.1/24
   ip address virtual 10.151.172.1/24 secondary
!
interface Vlan2104
   description DEV_DB1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.150.1/24
!
interface Vlan2105
   description DEV_APP1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.90.1/23
!
interface Vlan2106
   description Dev_DCC_DB1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.171.1/24
!
interface Vlan2107
   description DEV_DCP_DB1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.100.1/24
!
interface Vlan2108
   description DEV_DB1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.151.1/24
!
interface Vlan2109
   description SI_DCC_APP1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.186.1/23
!
interface Vlan2110
   description SI_DCC_DB2
   shutdown
   vrf DL_SI
   ip address virtual 10.151.185.1/24
!
interface Vlan2111
   description SI_Exadata
   shutdown
   vrf DL_SI
   ip address virtual 10.151.134.1/24
!
interface Vlan2112
   description SI_DCC_DB1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.181.1/24
!
interface Vlan2113
   description SI_DB1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.153.1/24
!
interface Vlan2114
   description SI_DCP_APP1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.191.1/24
!
interface Vlan2115
   description SI_APP1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.180.1/24
   ip address virtual 10.151.182.1/24 secondary
   ip address virtual 10.151.183.1/24 secondary
   ip address virtual 10.151.184.1/24 secondary
!
interface Vlan2116
   description SI_DCP_DB1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.192.1/24
!
interface Vlan2117
   description SI_APP1
   shutdown
   vrf DL_SI
   ip address virtual 10.151.152.1/24
!
interface Vlan2118
   description SI_NETAPP
   shutdown
   vrf DL_SI
   ip address virtual 10.151.255.1/24
!
interface Vlan2217
   description SAP-NIM
   no shutdown
   vrf DL_Prod
   ip address virtual 10.157.44.1/27
!
interface Vxlan1
   description ATL-ADM-BL201_VTEP
   vxlan source-interface Loopback1
   vxlan udp-port 4789
   vxlan vlan 21 vni 1021
   vxlan vlan 22 vni 10022
   vxlan vlan 25 vni 1025
   vxlan vlan 26 vni 10026
   vxlan vlan 1101 vni 11101
   vxlan vlan 1102 vni 11102
   vxlan vlan 1103 vni 11103
   vxlan vlan 1105 vni 11105
   vxlan vlan 1106 vni 11106
   vxlan vlan 1107 vni 11107
   vxlan vlan 1108 vni 11108
   vxlan vlan 1109 vni 11109
   vxlan vlan 1301 vni 11301
   vxlan vlan 1302 vni 11302
   vxlan vlan 1303 vni 11303
   vxlan vlan 1701 vni 11701
   vxlan vlan 1702 vni 11702
   vxlan vlan 1703 vni 11703
   vxlan vlan 1704 vni 11704
   vxlan vlan 1705 vni 11705
   vxlan vlan 1749 vni 11749
   vxlan vlan 1801 vni 11801
   vxlan vlan 1802 vni 11802
   vxlan vlan 1803 vni 11803
   vxlan vlan 2001 vni 12001
   vxlan vlan 2002 vni 12002
   vxlan vlan 2003 vni 12003
   vxlan vlan 2004 vni 12004
   vxlan vlan 2005 vni 12005
   vxlan vlan 2006 vni 12006
   vxlan vlan 2007 vni 12007
   vxlan vlan 2008 vni 12008
   vxlan vlan 2009 vni 12009
   vxlan vlan 2010 vni 12010
   vxlan vlan 2101 vni 12101
   vxlan vlan 2102 vni 12102
   vxlan vlan 2103 vni 12103
   vxlan vlan 2104 vni 12104
   vxlan vlan 2105 vni 12105
   vxlan vlan 2106 vni 12106
   vxlan vlan 2107 vni 12107
   vxlan vlan 2108 vni 12108
   vxlan vlan 2109 vni 12109
   vxlan vlan 2110 vni 12110
   vxlan vlan 2111 vni 12111
   vxlan vlan 2112 vni 12112
   vxlan vlan 2113 vni 12113
   vxlan vlan 2114 vni 12114
   vxlan vlan 2115 vni 12115
   vxlan vlan 2116 vni 12116
   vxlan vlan 2117 vni 12117
   vxlan vlan 2118 vni 12118
   vxlan vlan 2217 vni 12217
   vxlan vrf CORE vni 900
   vxlan vrf DL_Prod vni 901
   vxlan vrf DL_SI vni 904
!
ip virtual-router mac-address 00:1c:73:00:00:99
!
ip routing
ip routing vrf CORE
ip routing vrf DL_Prod
ip routing vrf DL_SI
no ip routing vrf MGMT
!
monitor session span_1 source Ethernet1/1-3/1 both
monitor session span_1 source Ethernet5/1-36/1 both
monitor session span_1 destination Ethernet4/1
!
ip community-list COM-ADM permit 65400:65400
ip community-list COM-ALPH permit 64600:64600
!
ip prefix-list ADM-PREFIX
   seq 10 permit 10.157.0.0/16 le 32
   seq 20 permit 10.151.0.0/16 le 32
!
ip prefix-list default
   seq 10 permit 0.0.0.0/0
!
ip prefix-list host-routes
   seq 10 permit 0.0.0.0/0 ge 32
!
ip prefix-list PL-LOOPBACKS-EVPN-OVERLAY
   seq 10 permit 10.157.4.0/23 eq 32
   seq 20 permit 10.157.6.0/23 eq 32
!
ip route vrf MGMT 0.0.0.0/0 10.157.18.1
ip route vrf CORE 10.157.0.0/16 null0
!
route-map allowed-default-only permit 10
   description allow default
   match ip address prefix-list default
   set community 65108:65108
!
route-map CL-ADM permit 10
   description adm community 65400:65400
   set community 65400:65400
   set weight 200
!
route-map CL-DEFAULT permit 10
   description adm community for default route 65108:65108
   set community 65108:65108
!
route-map IN-FM-CORE deny 10
   description adm community 65400:65400
   match community COM-ADM
!
route-map IN-FM-CORE permit 20
   description any
!
route-map OUT-TO-CORE permit 10
   description adm community 65400:65400
   match community COM-ADM
   match ip address prefix-list ADM-PREFIX
   set as-path match all replacement none
!
route-map remove-as-in permit 10
   description remove-as
   set as-path match all replacement auto
!
route-map remove-as-out deny 10
   description dont allow default
   match ip address prefix-list default
!
route-map remove-as-out deny 20
   description dont allow route type 2 only allow prefix route or route type 5
   match ip address prefix-list host-routes
!
route-map remove-as-out permit 30
   description remove-as for firewall to accept routes on same leaf
   set as-path match all replacement auto
!
route-map RM-CONN-2-BGP permit 10
   match ip address prefix-list PL-LOOPBACKS-EVPN-OVERLAY
!
route-map RM-VRF-CONN-2-BGP-CM permit 10
   description community 65400:65400
   set community 65400:65400
!
route-map sns-cm permit 10
   description community 65400:65400
   set community 65400:65400
!
router bfd
   multihop interval 300 min-rx 300 multiplier 3
!
router bgp 65401
   router-id 10.157.4.11
   distance bgp 20 200 200
   graceful-restart restart-time 300
   graceful-restart
   maximum-paths 4 ecmp 4
   update wait-install
   no bgp default ipv4-unicast
   neighbor ATL-ADM-CORE peer group
   neighbor ATL-ADM-CORE remote-as 65108
   neighbor ATL-ADM-CORE local-as 65499 no-prepend replace-as
   neighbor ATL-ADM-CORE description ATL-ADM-CORE
   neighbor ATL-ADM-CORE bfd
   neighbor ATL-ADM-CORE password 7 rY5j97zrotoM93IIdNJi0g==
   neighbor ATL-ADM-CORE send-community standard
   neighbor ATL-ADM-CORE maximum-routes 0
   neighbor ATL-ADM-FW-IN peer group
   neighbor ATL-ADM-FW-IN remote-as 65498
   neighbor ATL-ADM-FW-IN local-as 65497 no-prepend replace-as
   neighbor ATL-ADM-FW-IN update-source vlan 22
   neighbor ATL-ADM-FW-IN description ATL-ADM-FW-IN
   neighbor ATL-ADM-FW-IN password 7 tP/Gl0g5m6rTvt1fBbVpCg==
   neighbor ATL-ADM-FW-IN send-community standard
   neighbor ATL-ADM-FW-IN maximum-routes 0
   neighbor ATL-ADM-FW-OUT peer group
   neighbor ATL-ADM-FW-OUT remote-as 65498
   neighbor ATL-ADM-FW-OUT local-as 65499 no-prepend replace-as
   neighbor ATL-ADM-FW-OUT update-source vlan 21
   neighbor ATL-ADM-FW-OUT description ATL-ADM-FW-OUT
   neighbor ATL-ADM-FW-OUT password 7 Vb5Y7Hb+qvE2AePFRwuNaw==
   neighbor ATL-ADM-FW-OUT default-originate route-map CL-DEFAULT always
   neighbor ATL-ADM-FW-OUT send-community standard
   neighbor ATL-ADM-FW-OUT maximum-routes 0
   neighbor ATL-ADM-SI-FW-IN peer group
   neighbor ATL-ADM-SI-FW-IN remote-as 65496
   neighbor ATL-ADM-SI-FW-IN local-as 65495 no-prepend replace-as
   neighbor ATL-ADM-SI-FW-IN update-source vlan 26
   neighbor ATL-ADM-SI-FW-IN description ATL-ADM-SI-FW-IN
   neighbor ATL-ADM-SI-FW-IN password 7 t6GEeC9Y3FTCD0cDSmaZ+w==
   neighbor ATL-ADM-SI-FW-IN send-community
   neighbor ATL-ADM-SI-FW-IN maximum-routes 0
   neighbor ATL-ADM-SI-FW-OUT peer group
   neighbor ATL-ADM-SI-FW-OUT remote-as 65496
   neighbor ATL-ADM-SI-FW-OUT local-as 65499 no-prepend replace-as
   neighbor ATL-ADM-SI-FW-OUT update-source vlan 25
   neighbor ATL-ADM-SI-FW-OUT description ATL-ADM-SI-FW-OUT
   neighbor ATL-ADM-SI-FW-OUT password 7 sJLv8cLr9SLISJSVBnOWlg==
   neighbor ATL-ADM-SI-FW-OUT default-originate route-map CL-DEFAULT always
   neighbor ATL-ADM-SI-FW-OUT send-community standard
   neighbor ATL-ADM-SI-FW-OUT maximum-routes 0
   neighbor EVPN-OVERLAY-CORE peer group
   neighbor EVPN-OVERLAY-CORE update-source Loopback0
   neighbor EVPN-OVERLAY-CORE bfd
   neighbor EVPN-OVERLAY-CORE ebgp-multihop 15
   neighbor EVPN-OVERLAY-CORE password 7 rkymOvxfBRwlmqtxW6IjcQ==
   neighbor EVPN-OVERLAY-CORE send-community
   neighbor EVPN-OVERLAY-CORE maximum-routes 0
   neighbor EVPN-OVERLAY-PEERS peer group
   neighbor EVPN-OVERLAY-PEERS update-source Loopback0
   neighbor EVPN-OVERLAY-PEERS bfd
   neighbor EVPN-OVERLAY-PEERS ebgp-multihop 3
   neighbor EVPN-OVERLAY-PEERS password 7 FZCbOg206D07yEuxyTxVIg==
   neighbor EVPN-OVERLAY-PEERS send-community
   neighbor EVPN-OVERLAY-PEERS maximum-routes 0
   neighbor IPv4-UNDERLAY-PEERS peer group
   neighbor IPv4-UNDERLAY-PEERS password 7 YCwcJual/+vy4x/M32Hu5A==
   neighbor IPv4-UNDERLAY-PEERS send-community
   neighbor IPv4-UNDERLAY-PEERS maximum-routes 12000
   neighbor 10.157.0.60 peer group IPv4-UNDERLAY-PEERS
   neighbor 10.157.0.60 remote-as 65400
   neighbor 10.157.0.60 description ATL-ADM-SP101_Ethernet1/1
   neighbor 10.157.0.62 peer group IPv4-UNDERLAY-PEERS
   neighbor 10.157.0.62 remote-as 65400
   neighbor 10.157.0.62 description ATL-ADM-SP102_Ethernet1/1
   neighbor 10.157.0.64 peer group IPv4-UNDERLAY-PEERS
   neighbor 10.157.0.64 remote-as 65400
   neighbor 10.157.0.64 description ATL-ADM-SP103_Ethernet1/1
   neighbor 10.157.4.1 peer group EVPN-OVERLAY-PEERS
   neighbor 10.157.4.1 remote-as 65400
   neighbor 10.157.4.1 description ATL-ADM-SP101
   neighbor 10.157.4.2 peer group EVPN-OVERLAY-PEERS
   neighbor 10.157.4.2 remote-as 65400
   neighbor 10.157.4.2 description ATL-ADM-SP102
   neighbor 10.157.4.3 peer group EVPN-OVERLAY-PEERS
   neighbor 10.157.4.3 remote-as 65400
   neighbor 10.157.4.3 description ATL-ADM-SP103
   redistribute connected route-map RM-CONN-2-BGP
   !
   vlan 1101
      rd 10.157.4.11:11101
      rd evpn domain remote 10.157.4.11:11101
      route-target both 11101:11101
      route-target import export evpn domain remote 11101:11101
      redistribute learned
   !
   vlan 1102
      rd 10.157.4.11:11102
      rd evpn domain remote 10.157.4.11:11102
      route-target both 11102:11102
      route-target import export evpn domain remote 11102:11102
      redistribute learned
   !
   vlan 1103
      rd 10.157.4.11:11103
      rd evpn domain remote 10.157.4.11:11103
      route-target both 11103:11103
      route-target import export evpn domain remote 11103:11103
      redistribute learned
   !
   vlan 1105
      rd 10.157.4.11:11105
      rd evpn domain remote 10.157.4.11:11105
      route-target both 11105:11105
      route-target import export evpn domain remote 11105:11105
      redistribute learned
   !
   vlan 1106
      rd 10.157.4.11:11106
      rd evpn domain remote 10.157.4.11:11106
      route-target both 11106:11106
      route-target import export evpn domain remote 11106:11106
      redistribute learned
   !
   vlan 1107
      rd 10.157.4.11:11107
      rd evpn domain remote 10.157.4.11:11107
      route-target both 11107:11107
      route-target import export evpn domain remote 11107:11107
      redistribute learned
   !
   vlan 1108
      rd 10.157.4.11:11108
      rd evpn domain remote 10.157.4.11:11108
      route-target both 11108:11108
      route-target import export evpn domain remote 11108:11108
      redistribute learned
   !
   vlan 1109
      rd 10.157.4.11:11109
      rd evpn domain remote 10.157.4.11:11109
      route-target both 11109:11109
      route-target import export evpn domain remote 11109:11109
      redistribute learned
   !
   vlan 1301
      rd 10.157.4.11:11301
      rd evpn domain remote 10.157.4.11:11301
      route-target both 11301:11301
      route-target import export evpn domain remote 11301:11301
      redistribute learned
   !
   vlan 1302
      rd 10.157.4.11:11302
      rd evpn domain remote 10.157.4.11:11302
      route-target both 11302:11302
      route-target import export evpn domain remote 11302:11302
      redistribute learned
   !
   vlan 1303
      rd 10.157.4.11:11303
      rd evpn domain remote 10.157.4.11:11303
      route-target both 11303:11303
      route-target import export evpn domain remote 11303:11303
      redistribute learned
   !
   vlan 1701
      rd 10.157.4.11:11701
      rd evpn domain remote 10.157.4.11:11701
      route-target both 11701:11701
      route-target import export evpn domain remote 11701:11701
      redistribute learned
   !
   vlan 1702
      rd 10.157.4.11:11702
      rd evpn domain remote 10.157.4.11:11702
      route-target both 11702:11702
      route-target import export evpn domain remote 11702:11702
      redistribute learned
   !
   vlan 1703
      rd 10.157.4.11:11703
      rd evpn domain remote 10.157.4.11:11703
      route-target both 11703:11703
      route-target import export evpn domain remote 11703:11703
      redistribute learned
   !
   vlan 1704
      rd 10.157.4.11:11704
      rd evpn domain remote 10.157.4.11:11704
      route-target both 11704:11704
      route-target import export evpn domain remote 11704:11704
      redistribute learned
   !
   vlan 1705
      rd 10.157.4.11:11705
      rd evpn domain remote 10.157.4.11:11705
      route-target both 11705:11705
      route-target import export evpn domain remote 11705:11705
      redistribute learned
   !
   vlan 1749
      rd 10.157.4.11:11749
      rd evpn domain remote 10.157.4.11:11749
      route-target both 11749:11749
      route-target import export evpn domain remote 11749:11749
      redistribute learned
   !
   vlan 1801
      rd 10.157.4.11:11801
      rd evpn domain remote 10.157.4.11:11801
      route-target both 11801:11801
      route-target import export evpn domain remote 11801:11801
      redistribute learned
   !
   vlan 1802
      rd 10.157.4.11:11802
      rd evpn domain remote 10.157.4.11:11802
      route-target both 11802:11802
      route-target import export evpn domain remote 11802:11802
      redistribute learned
   !
   vlan 1803
      rd 10.157.4.11:11803
      rd evpn domain remote 10.157.4.11:11803
      route-target both 11803:11803
      route-target import export evpn domain remote 11803:11803
      redistribute learned
   !
   vlan 2001
      rd 10.157.4.11:12001
      rd evpn domain remote 10.157.4.11:12001
      route-target both 12001:12001
      route-target import export evpn domain remote 12001:12001
      redistribute learned
   !
   vlan 2002
      rd 10.157.4.11:12002
      rd evpn domain remote 10.157.4.11:12002
      route-target both 12002:12002
      route-target import export evpn domain remote 12002:12002
      redistribute learned
   !
   vlan 2003
      rd 10.157.4.11:12003
      rd evpn domain remote 10.157.4.11:12003
      route-target both 12003:12003
      route-target import export evpn domain remote 12003:12003
      redistribute learned
   !
   vlan 2004
      rd 10.157.4.11:12004
      rd evpn domain remote 10.157.4.11:12004
      route-target both 12004:12004
      route-target import export evpn domain remote 12004:12004
      redistribute learned
   !
   vlan 2005
      rd 10.157.4.11:12005
      rd evpn domain remote 10.157.4.11:12005
      route-target both 12005:12005
      route-target import export evpn domain remote 12005:12005
      redistribute learned
   !
   vlan 2006
      rd 10.157.4.11:12006
      rd evpn domain remote 10.157.4.11:12006
      route-target both 12006:12006
      route-target import export evpn domain remote 12006:12006
      redistribute learned
   !
   vlan 2007
      rd 10.157.4.11:12007
      rd evpn domain remote 10.157.4.11:12007
      route-target both 12007:12007
      route-target import export evpn domain remote 12007:12007
      redistribute learned
   !
   vlan 2008
      rd 10.157.4.11:12008
      rd evpn domain remote 10.157.4.11:12008
      route-target both 12008:12008
      route-target import export evpn domain remote 12008:12008
      redistribute learned
   !
   vlan 2009
      rd 10.157.4.11:12009
      rd evpn domain remote 10.157.4.11:12009
      route-target both 12009:12009
      route-target import export evpn domain remote 12009:12009
      redistribute learned
   !
   vlan 2010
      rd 10.157.4.11:12010
      rd evpn domain remote 10.157.4.11:12010
      route-target both 12010:12010
      route-target import export evpn domain remote 12010:12010
      redistribute learned
   !
   vlan 21
      rd 10.157.4.11:1021
      rd evpn domain remote 10.157.4.11:1021
      route-target both 1021:1021
      route-target import export evpn domain remote 1021:1021
      redistribute learned
   !
   vlan 2101
      rd 10.157.4.11:12101
      rd evpn domain remote 10.157.4.11:12101
      route-target both 12101:12101
      route-target import export evpn domain remote 12101:12101
      redistribute learned
   !
   vlan 2102
      rd 10.157.4.11:12102
      rd evpn domain remote 10.157.4.11:12102
      route-target both 12102:12102
      route-target import export evpn domain remote 12102:12102
      redistribute learned
   !
   vlan 2103
      rd 10.157.4.11:12103
      rd evpn domain remote 10.157.4.11:12103
      route-target both 12103:12103
      route-target import export evpn domain remote 12103:12103
      redistribute learned
   !
   vlan 2104
      rd 10.157.4.11:12104
      rd evpn domain remote 10.157.4.11:12104
      route-target both 12104:12104
      route-target import export evpn domain remote 12104:12104
      redistribute learned
   !
   vlan 2105
      rd 10.157.4.11:12105
      rd evpn domain remote 10.157.4.11:12105
      route-target both 12105:12105
      route-target import export evpn domain remote 12105:12105
      redistribute learned
   !
   vlan 2106
      rd 10.157.4.11:12106
      rd evpn domain remote 10.157.4.11:12106
      route-target both 12106:12106
      route-target import export evpn domain remote 12106:12106
      redistribute learned
   !
   vlan 2107
      rd 10.157.4.11:12107
      rd evpn domain remote 10.157.4.11:12107
      route-target both 12107:12107
      route-target import export evpn domain remote 12107:12107
      redistribute learned
   !
   vlan 2108
      rd 10.157.4.11:12108
      rd evpn domain remote 10.157.4.11:12108
      route-target both 12108:12108
      route-target import export evpn domain remote 12108:12108
      redistribute learned
   !
   vlan 2109
      rd 10.157.4.11:12109
      rd evpn domain remote 10.157.4.11:12109
      route-target both 12109:12109
      route-target import export evpn domain remote 12109:12109
      redistribute learned
   !
   vlan 2110
      rd 10.157.4.11:12110
      rd evpn domain remote 10.157.4.11:12110
      route-target both 12110:12110
      route-target import export evpn domain remote 12110:12110
      redistribute learned
   !
   vlan 2111
      rd 10.157.4.11:12111
      rd evpn domain remote 10.157.4.11:12111
      route-target both 12111:12111
      route-target import export evpn domain remote 12111:12111
      redistribute learned
   !
   vlan 2112
      rd 10.157.4.11:12112
      rd evpn domain remote 10.157.4.11:12112
      route-target both 12112:12112
      route-target import export evpn domain remote 12112:12112
      redistribute learned
   !
   vlan 2113
      rd 10.157.4.11:12113
      rd evpn domain remote 10.157.4.11:12113
      route-target both 12113:12113
      route-target import export evpn domain remote 12113:12113
      redistribute learned
   !
   vlan 2114
      rd 10.157.4.11:12114
      rd evpn domain remote 10.157.4.11:12114
      route-target both 12114:12114
      route-target import export evpn domain remote 12114:12114
      redistribute learned
   !
   vlan 2115
      rd 10.157.4.11:12115
      rd evpn domain remote 10.157.4.11:12115
      route-target both 12115:12115
      route-target import export evpn domain remote 12115:12115
      redistribute learned
   !
   vlan 2116
      rd 10.157.4.11:12116
      rd evpn domain remote 10.157.4.11:12116
      route-target both 12116:12116
      route-target import export evpn domain remote 12116:12116
      redistribute learned
   !
   vlan 2117
      rd 10.157.4.11:12117
      rd evpn domain remote 10.157.4.11:12117
      route-target both 12117:12117
      route-target import export evpn domain remote 12117:12117
      redistribute learned
   !
   vlan 2118
      rd 10.157.4.11:12118
      rd evpn domain remote 10.157.4.11:12118
      route-target both 12118:12118
      route-target import export evpn domain remote 12118:12118
      redistribute learned
   !
   vlan 22
      rd 10.157.4.11:10022
      rd evpn domain remote 10.157.4.11:10022
      route-target both 10022:10022
      route-target import export evpn domain remote 10022:10022
      redistribute learned
   !
   vlan 2217
      rd 10.157.4.11:12217
      rd evpn domain remote 10.157.4.11:12217
      route-target both 12217:12217
      route-target import export evpn domain remote 12217:12217
      redistribute learned
   !
   vlan 25
      rd 10.157.4.11:1025
      rd evpn domain remote 10.157.4.11:1025
      route-target both 1025:1025
      route-target import export evpn domain remote 1025:1025
      redistribute learned
   !
   vlan 26
      rd 10.157.4.11:10026
      rd evpn domain remote 10.157.4.11:10026
      route-target both 10026:10026
      route-target import export evpn domain remote 10026:10026
      redistribute learned
   !
   address-family evpn
      neighbor EVPN-OVERLAY-CORE activate
      neighbor EVPN-OVERLAY-CORE domain remote
      neighbor EVPN-OVERLAY-PEERS activate
      neighbor default next-hop-self received-evpn-routes route-type ip-prefix inter-domain
   !
   address-family ipv4
      no neighbor EVPN-OVERLAY-CORE activate
      no neighbor EVPN-OVERLAY-PEERS activate
      neighbor IPv4-UNDERLAY-PEERS activate
   !
   vrf CORE
      rd 10.157.4.11:900
      route-target import evpn 900:900
      route-target export evpn 900:900
      router-id 10.157.4.11
      update wait-install
      neighbor 10.157.13.1 peer group ATL-ADM-CORE
      neighbor 10.157.13.1 bfd
      neighbor 10.157.13.1 route-map OUT-TO-CORE out
      neighbor 10.157.13.1 route-map IN-FM-CORE in
      neighbor 10.157.13.3 peer group ATL-ADM-CORE
      neighbor 10.157.13.3 bfd
      neighbor 10.157.13.3 route-map OUT-TO-CORE out
      neighbor 10.157.13.3 route-map IN-FM-CORE in
      neighbor 10.157.14.1 peer group ATL-ADM-FW-OUT
      neighbor 10.157.14.1 bfd
      neighbor 10.157.14.1 timers 5 15
      neighbor 10.157.14.1 route-map allowed-default-only out
      neighbor 10.157.14.33 peer group ATL-ADM-SI-FW-OUT
      neighbor 10.157.14.33 bfd
      neighbor 10.157.14.33 timers 5 15
      neighbor 10.157.14.33 route-map allowed-default-only out
      network 10.157.0.0/16 route-map CL-ADM
      redistribute connected
      !
      address-family ipv4
         neighbor 10.157.13.1 activate
         neighbor 10.157.13.3 activate
         neighbor 10.157.14.1 activate
         neighbor 10.157.14.33 activate
   !
   vrf DL_Prod
      rd 10.157.4.11:901
      route-target import evpn 901:901
      route-target export evpn 901:901
      router-id 10.157.4.11
      update wait-install
      neighbor 10.157.14.9 peer group ATL-ADM-FW-IN
      neighbor 10.157.14.9 bfd
      neighbor 10.157.14.9 timers 5 15
      neighbor 10.157.14.9 route-map remove-as-out out
      neighbor 10.157.14.9 route-map remove-as-in in
      redistribute connected route-map RM-VRF-CONN-2-BGP-CM
      !
      address-family ipv4
         neighbor 10.157.14.9 activate
   !
   vrf DL_SI
      rd 10.157.4.11:904
      route-target import evpn 904:904
      route-target export evpn 904:904
      router-id 10.157.4.11
      update wait-install
      neighbor 10.157.14.41 peer group ATL-ADM-SI-FW-IN
      neighbor 10.157.14.41 bfd
      neighbor 10.157.14.41 timers 5 15
      neighbor 10.157.14.41 route-map remove-as-out out
      neighbor 10.157.14.41 route-map remove-as-in in
      redistribute connected route-map RM-VRF-CONN-2-BGP-CM
      !
      address-family ipv4
         neighbor 10.157.14.41 activate
!
ip tacacs vrf MGMT source-interface Vlan201
!
banner login
&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&
&&&&     Access is restricted to     &&&&
&&&&      authorized users only.     &&&&
&&&&                                 &&&&
&&&&    Unauthorized access is a     &&&&
&&&& violation of state and federal, &&&&
&&&&     civil and criminal laws!    &&&&
&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&
EOF

!
banner motd
&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&
&&&&             Delta               &&&&
&&&&       Network Technology        &&&&
&&&&                                 &&&&
&&&&  -- Unauthorized Access is --   &&&&
&&&&  !!! STRICTLY PROHIBITED !!!    &&&&
&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&
EOF

!
management api http-commands
   protocol https
   no shutdown
   !
   vrf MGMT
      no shutdown
!
management console
   idle-timeout 10
!
end


`
