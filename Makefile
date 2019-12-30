.PHONY: fmt
fmt:
	gofmt -w *.go

.PHONY: build
build:
	go build .
