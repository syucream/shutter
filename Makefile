.PHONY: fmt
fmt:
	gofmt -w *.go
	gofmt -w cmd/**/*.go

.PHONY: build
build:
	go build cmd/shutter/shutter.go

.PHONY: test
test:
	go test -cover .
