# CRISPY ENIGMA

This project is a demo to use the artifacts created by [Arista AVD](https://avd.sh) during the build process and update other systems.
When the build playbook of AVD runs creates a folder called `intended\structure_configs` with a yml file for every switch in the topology. This project uses this files to update [Netbox](https://netbox.readthedocs.io/en/stable/) and [Infoblox](https://www.infoblox.com/) with the new information. This script also is capable of generating YAML files to be consumed by the `ansible-playbook build.yml` from AVD.

## Installation

- Download the correct version of your OS from the [releases](https://github.com/rogerscuall/crispy-enigma/releases) page.
- Build the binary(recommended):
  - Install [Go](https://golang.org/doc/install).
  - Download the repo: `git clone https://github.com/rogerscuall/crispy-enigma.git`.
  - Move to the repo: `cd crispy-enigma`.
  - Create the binary: `go build -o crispy-enigma`.

## How to use it

1. Run the AVD build playbook. `ansible-playbook -i build.yml`
1. Verify that the `intended\structure_configs` folder was created.
1. Create the necessary envar:
   1. `export INFOBLOX_USERNAME=<user>`
   1. `export INFOBLOX_PASSWORD=<password>`
   1. `export INFOBLOX_WAPI_VERSION=<version>`
   1. `export INFOBLOX_URL=<url>` -> `export INFOBLOX_URL=1.1.1.1`
   1. `export NETBOX_USERNAME=<user>`
   1. `export NETBOX_PASSWORD=<password>`
   1. `export NETBOX_URL=<url>` -> `export NETBOX_URL=https://demo.netbox.dev`
1. Update Infoblox: `crispy-enigma infobloxUpdate -f intended/structure_configs/`
1. Update Netbox: `crispy-enigma netboxUpdate -f intended/structure_configs/`

## Known issues

- Currently there is a problem in the netbox client API [issue](https://github.com/netbox-community/go-netbox/issues/164) so is recommended to build using the vendor folder provided in the repo.
