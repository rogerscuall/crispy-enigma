package models

type Network struct {
	Configs []*Config `yaml:"configs"`
}

func (n Network) GetHostnames() []string {
	var hostnames []string
	for _, c := range n.Configs {
		hostnames = append(hostnames, c.Hostname)
	}
	return hostnames
}

type Links struct {
	Links []Connection `yaml:"connections"`
}

func (l Links) GetLinks() []Connection {
	return l.Links
}

type Connection struct {
	SideA string `yaml:"side_a"`
	SideB string `yaml:"side_b"`
	PortA string `yaml:"port_a"`
	PortB string `yaml:"port_b"`
}

/*
GetNetworkConnections returns the connections that are inside the network
*/
func (n Network) GetNetworkConnections() []Connection {
	//TODO: Seems like this function does the same as GetInNetworkConnections
	hostnames := n.GetHostnames()
	networkInterfaces := []Connection{}
	for _, hostname := range hostnames {
		for _, c := range n.Configs {
			if c.Hostname == hostname {
				for _, e := range c.EthernetInterfaces {
					if e.Peer == "" || e.Peer == "UNUSED" {
						continue
					}
					connection := Connection{
						SideA: hostname,
						SideB: e.Peer,
						PortA: e.Name,
						PortB: e.PeerInterface,
					}
					networkInterfaces = append(networkInterfaces, connection)
				}
			}
		}
	}
	return networkInterfaces
}

// GetInNetworkConnections returns the connections that are inside the network
// This are the connections that connect two devices in the network
func (n Network) GetInNetworkConnections() []Connection {
	hostnames := n.GetHostnames()
	networkInterfaces := []Connection{}
	for _, config := range n.Configs {
		for _, e := range config.EthernetInterfaces {
			for _, hostname := range hostnames {
				if e.Peer == hostname {
					connection := Connection{
						SideA: config.Hostname,
						SideB: e.Peer,
						PortA: e.Name,
						PortB: e.PeerInterface,
					}
					networkInterfaces = append(networkInterfaces, connection)
				}
			}
		}
	}
	return networkInterfaces
}

// CleanNetworkConnections guarantees that the connections are unique
func (n Network) CleanNetworkConnections() []Connection {
	connections := n.GetNetworkConnections()
	uniqueConnections := []Connection{}
	for _, c := range connections {
		if !containsConnection(uniqueConnections, c) {
			uniqueConnections = append(uniqueConnections, c)
		}
	}
	return uniqueConnections
}

/*
CleanInNetworkConnections guarantees that the connections are unique
*/
func (n Network) CleanInNetworkConnections() []Connection {
	connections := n.GetInNetworkConnections()
	uniqueConnections := []Connection{}
	for _, c := range connections {
		if !containsConnection(uniqueConnections, c) {
			uniqueConnections = append(uniqueConnections, c)
		}
	}
	return uniqueConnections
}

func containsConnection(connections []Connection, connection Connection) bool {
	for _, c := range connections {
		if c.SideA == connection.SideA && c.SideB == connection.SideB {
			return true
		}
		if c.SideA == connection.SideB && c.SideB == connection.SideA {
			return true
		}
	}
	return false
}
