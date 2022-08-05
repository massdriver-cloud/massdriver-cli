INSTALL_PATH ?= /usr/local/bin
GIT_SHA := $(shell git log -1 --pretty=format:"%H")
LD_FLAGS := "-X github.com/massdriver-cloud/massdriver-cli/pkg/version.version=dev -X github.com/massdriver-cloud/massdriver-cli/pkg/version.gitSHA=local-dev-${GIT_SHA}"

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

.PHONY refresh-templates
refresh-templates:
	./scripts/refresh-templates.sh
