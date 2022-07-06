INSTALL_PATH ?= /usr/local/bin

.PHONY: test
test:
	go test ./cmd
	go test ./pkg/...
	go build -o ./mass

bin:
	mkdir bin

.PHONY: build
build: bin
	GOOS=darwin GOARCH=arm64 go build -o bin/mass-darwin-arm64
	GOOS=linux GOARCH=amd64 go build -o bin/mass-linux-amd64

.PHONY: install.macos
install.macos: build
	rm -f ${INSTALL_PATH}/mass
	cp bin/mass-darwin-arm64 ${INSTALL_PATH}/mass

.PHONY: install.linux
install.linux: build
	cp -f bin/mass-linux-amd64 ${INSTALL_PATH}/mass
