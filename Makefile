.PHONY: test fmt vet build

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

build:
	go build ./cmd/ffb
