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
