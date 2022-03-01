.PHONY: test
test:
	go test ./cmd
	go test ./src/...
	go build -o ./mass

build:
	GOOS=linux GOARCH=amd64 go build -o ./mass
