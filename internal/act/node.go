package act

type Node struct {
	Name              string   `yaml:"-"`
	IPAddr            string   `yaml:"ip_addr"`
	NodeType          string   `yaml:"node_type"`
	AutoConfiguration bool     `yaml:"auto_configuration,omitempty"`
	Ports             []string `yaml:"ports,omitempty"`
}

func (n *Node) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var nodeMap map[string]interface{}
	if err := unmarshal(&nodeMap); err != nil {
		return err
	}

	for key, value := range nodeMap {
		n.Name = key

		nodeDetails := value.(map[interface{}]interface{})
		n.IPAddr = nodeDetails["ip_addr"].(string)
		n.NodeType = nodeDetails["node_type"].(string)
		if autoConfig, ok := nodeDetails["auto_configuration"].(bool); ok {
			n.AutoConfiguration = autoConfig
		}
		if ports, ok := nodeDetails["ports"].([]interface{}); ok {
			for _, port := range ports {
				portName := port.(string)
				n.Ports = append(n.Ports, portName)
			}
		}
	}

	return nil
}

func (n Node) MarshalYAML() (interface{}, error) {
	nodeDetails := map[string]interface{}{}

	// Only add non-empty fields to the nodeDetails map
	if n.IPAddr != "" {
		nodeDetails["ip_addr"] = n.IPAddr
	}
	if n.NodeType != "" {
		nodeDetails["node_type"] = n.NodeType
	}
	if n.AutoConfiguration {
		nodeDetails["auto_configuration"] = n.AutoConfiguration
	}
	if len(n.Ports) > 0 {
		nodeDetails["ports"] = n.Ports
	}

	// Return a map with the node name as the key, only if nodeDetails is not empty
	if len(nodeDetails) > 0 {
		return map[string]interface{}{
			n.Name: nodeDetails,
		}, nil
	}

	return nil, nil // Avoid including an empty map
}
