.PHONY: test
test:
	go test ./cmd
	go test ./pkg/...
	go build -o ./mass

build:
	GOOS=linux GOARCH=amd64 go build -o ./mass

build-m1:
	GOOS=darwin GOARCH=arm64 go build -o ./mass
