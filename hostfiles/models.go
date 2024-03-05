package hostfiles

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Interface struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Shutdown    bool   `yaml:"shutdown"`
	VLAN        int `yaml:"vlan,omitempty"`
}

func NewInterfaces(initial, end int) []Interface {
	var interfaces []Interface
	for initial <= end {
		interfaces = append(interfaces, Interface{
			Shutdown: true,
		})
		initial++
	}
	return interfaces
}

func WriteYamlFile(file string, interfaces []Interface) {
	// Create the file
	f, err := os.Create(strings.Replace(file, ".csv", ".yml", 1))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Write the interfaces to the file
	var config struct {
		CscEthernetInterfaces []Interface `yaml:"custom_structured_configuration_ethernet_interfaces"`
	}
	config.CscEthernetInterfaces = interfaces
	err = yaml.NewEncoder(f).Encode(config)
	if err != nil {
		panic(err)
	}
}

// WriteHostFile writes the interfaces to a file
func WriteHostFile(file string, interfaces []Interface) {
	// Create the file
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Write the interfaces to the file
	var config struct {
		CscEthernetInterfaces []Interface `yaml:"custom_structured_configuration_ethernet_interfaces"`
	}
	config.CscEthernetInterfaces = interfaces
	err = yaml.NewEncoder(f).Encode(config)
	if err != nil {
		panic(err)
	}
}

// ParseInterfaceString parses the string and returns the interface name format, lower and higher bounds.
// Currently only parses the format EthernetX-Y
func ParseInterfaceString(s string) (string, int, int, error) {
	// TODO: Parse the different options like Ethernet1/1
	re := regexp.MustCompile(`([a-zA-Z\/]+)(\d+)-(\d+)`)
	matches := re.FindStringSubmatch(s)

	if len(matches) < 4 {
		return "", 0, 0, fmt.Errorf("invalid format")
	}

	base := matches[1]
	lower, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid lower bound")
	}

	higher, err := strconv.Atoi(matches[3])
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid higher bound")
	}

	return base, lower, higher, nil
}

// CreateDefaultInterfaces creates a default set of interfaces
func CreateDefaultInterfaces(base string, lower, higher int) []Interface {
	var interfaces []Interface
	for lower <= higher {
		interfaces = append(interfaces, Interface{
			Name:        fmt.Sprintf("%s%d", base, lower),
			Description: "unused",
			Shutdown:    true,
		})
		lower++
	}
	return interfaces
}
