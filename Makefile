include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/infoblox/update: run the cmd/api application
.PHONY: run/infoblox/update
run/infoblox/update:
	go run *.go infobloxUpdate -f old

## run/infoblox/nextnetwork: run the infoblox next network process
.PHONY: run/infoblox/nextnetwork
run/infoblox/nextnetwork:
	go run *.go infobloxNextNetwork -f old/NETWORK_SERVICES.yml

## run/netbox/update: run the cmd/api application
.PHONY: run/netbox/update
run/netbox/update:
	go run *.go netboxUpdate -f old1

## run/host/interface: run the cmd/api application
.PHONY: run/host/interface
run/host/interface:
	go run *.go hostInterfaces -f old3

## run/cvp/config: run the cmd/api application
.PHONY: run/cvp/config
run/cvp/config:
	go run *.go cvpConfig -f old4

## run/cvp/config/check-all: run the cmd/api application
.PHONY: run/cvp/config/check-all
run/cvp/config/check-all:
	go run *.go cvpConfig -f old4 -c

## run/cvp/config/debug: run the cmd/api application
.PHONY: run/cvp/config/debug
run/cvp/config/debug:
	go run *.go cvpConfig -f old4 -v

## run/cvp/pending/tasks/debug: run the cmd/api application
.PHONY: run/cvp/pending/tasks/debug
run/cvp/pending/tasks/debug:
	go run *.go cvpPendingTask -f old4 -v

## run/pbAccessInterface: run the Pocket Base CMD
.PHONY: run/pbAccessInterface
run/pbAccessInterface:
	go run *.go pbAccessInterface -f old1 -v

## run/parseOutput: run the cmd/api application
.PHONY: run/parseOutput
run/parseOutput:
	go run *.go parseOutput -k warning,error -f old5/test.txt

## run/actTopology: run the actTopology CMD
.PHONY: run/actTopology
run/actTopology:
	go run *.go actTopology

## run/actInventory: run the actInventory CMD
.PHONY: run/actInventory
run/actInventory:
	go run *.go actInventory -i inventory.yml

## install: build and install the CLI locally
.PHONY: install
install:
	@echo 'Building and installing the CLI...'
	go build -o $(GOPATH)/bin/crispy-enigma
	@echo 'Installation complete. You can now use "crispy-enigma" command.'

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor


# ==================================================================================== #
# RELEASE
# ==================================================================================== #

## release: build and release the application
.PHONY: release
release:
	@echo 'Creating release...'
	goreleaser release --snapshot --clean

# ==================================================================================== #
# CONTAINER
# ==================================================================================== #

## build/container: build the container image
.PHONY: build/container
build/container:
	@echo 'Building container image...'
	podman build -t registry.presidio.com/arista/arista-avd-cvaas/crispy-enigma:${CE_TAG} .
	podman push registry.presidio.com/arista/arista-avd-cvaas/crispy-enigma:${CE_TAG}
