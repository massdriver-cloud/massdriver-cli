INSTALL_PATH ?= /usr/local/bin
GIT_SHA := $(shell git log -1 --pretty=format:"%H")
LD_FLAGS := "-X github.com/massdriver-cloud/massdriver-cli/pkg/version.version=dev -X github.com/massdriver-cloud/massdriver-cli/pkg/version.gitSHA=local-dev-${GIT_SHA}"

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))

.PHONY: test
test:
	go test ./cmd
	go test ./pkg/...

bin:
	mkdir bin

.PHONY: build
build: bin
	GOOS=darwin GOARCH=arm64 go build -o bin/mass-darwin-arm64 -ldflags=${LD_FLAGS}

.PHONY: build.linux
build.linux: bin
	GOOS=linux GOARCH=amd64 go build -o bin/mass-linux-amd64 -ldflags=${LD_FLAGS}

.PHONY: install
install: build
	rm -f ${INSTALL_PATH}/mass
	cp bin/mass-darwin-arm64 ${INSTALL_PATH}/mass

.PHONY: install.macos
install.macos: build
	rm -f ${INSTALL_PATH}/mass
	cp bin/mass-darwin-arm64 ${INSTALL_PATH}/mass


.PHONY: install.linux
install.linux: build.linux
	cp -f bin/mass-linux-amd64 ${INSTALL_PATH}/mass

# build the api2 client
MASSDRIVER_PATH?=../massdriver
API2_PATH?=${MKFILE_DIR}/pkg/api2
api:
	go get github.com/vektah/gqlparser/v2/validator@v2.5.1
	go get github.com/Khan/genqlient/generate@v0.5.0
	cd ${MASSDRIVER_PATH} && \
		mix absinthe.schema.sdl ${API2_PATH}/schema.graphql && \
		cd ${API2_PATH} && \
		go generate
