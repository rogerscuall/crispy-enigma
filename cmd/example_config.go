package cmd

var actExampleConfig = `
veos:
  password: cvp123!
  username: cvpadmin
  version: 4.27.0F

cvp:
  password: cvproot
  username: root
  version: 2022.2.2

generic:
  password: ansible
  username: ansible
  version: Rocky-8.5

third-party:
  password: ansible
  username: ansible
  version: infoblox

nodes:
# Update the IP addresses of the test nodes and CVP to be in the same subnet as the devices.
  - CVP:
      ip_addr: 192.168.0.5
      node_type: cvp
      auto_configuration: true
  - INTERNAL-TEST:
      ip_addr: 192.168.0.11
      node_type: veos
      ports:
        - Ethernet1-32
  - EXTERNAL-TEST:
      ip_addr: 192.168.0.11
      node_type: veos
      ports:
        - Ethernet1-32
links:
  - connection:
      - INTERNAL-TEST:Ethernet1
      # Update this next line with the name of an existing switch and a port that is configured but not connected to another switch in this network
      - DC1-L2LEAF1A:Ethernet20
  - connection:
      - EXTERNAL-TEST:Ethernet1
      # Update this next line with the name of an existing switch and a port that is configured but not connected to another switch in this network
      - DC1-L2LEAF1A:Ethernet21
`
