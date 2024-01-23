package hostfiles

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Interface struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Shutdown    bool   `yaml:"shutdown"`
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
