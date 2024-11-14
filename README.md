# CRISPY ENIGMA

This project is a demo to use the artifacts created by [Arista AVD](https://avd.sh).
Can compare the new configuration created with AVD and the running configuration in CVP.
It can automate the process of creating an ACT lab from an AVD topology.

## Installation

- (Recommended) Download the correct version of your OS from the [releases](https://github.com/rogerscuall/crispy-enigma/releases) page.
- Build the binary:
  - Install [Go](https://golang.org/doc/install).
  - Download the repo: `git clone https://github.com/rogerscuall/crispy-enigma.git`.
  - Move to the repo: `cd crispy-enigma`.
  - Create the binary: `go build -o crispy-enigma`.
- Use the docker image: `docker pull rogerscuall/crispy-enigma`.
- Download the binary from the releases.

## AVD CVP Compare

Download running-config from CVP and compares with intended config and reports the differences.

1. Run the AVD build playbook. `ansible-playbook -i build.yml`
1. Verify that the `intended\structure_configs` folder was created.
1. Create the necessary envar:
   1. `export CVP_URL=<url>` -> if not defined uses `https://www.arista.io`
   1. `export CVP_USERNAME=<user>`
   1. `export CVP_PASSWORD=<password>`
1. Download and compare configs with CVP: `crispy-enigma cvpConfig -f intended/configs/`

## ACT Automation

[ACT Automation](documentation/ACT/README.md)
