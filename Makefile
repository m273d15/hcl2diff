build:
	go mod tidy
	go build

test:
	go test -v ./...

fmt:
	go fmt ./...
