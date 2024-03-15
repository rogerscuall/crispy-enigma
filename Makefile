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

## run/cvp/config/debug: run the cmd/api application
.PHONY: run/cvp/config/debug
run/cvp/config/debug:
	go run *.go cvpConfig -f old4 -v


## run/parseOutput: run the cmd/api application
.PHONY: run/parseOutput
run/parseOutput:
	go run *.go parseOutput -k warning,error -f old5/test.txt

## db/migrations/new name=<VALUE>ssh-rsa 
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

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