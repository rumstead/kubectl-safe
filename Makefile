.PHONY: build
build:
	go build .

.PHONY: test
test:
	go test ./... -coverprofile coverage.out