package models

func GetHostsFromConfig(config []*Config) []string {
	var hostnames []string
	for _, c := range config {
		hostname := c.Hostname
		hostnames = append(hostnames, hostname)
	}
	return hostnames
}
