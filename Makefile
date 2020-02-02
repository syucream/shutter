.PHONY: init
init:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: fmt
fmt:
	gofmt -w *.go
	gofmt -w cmd/**/*.go

.PHONY: build
build:
	go build cmd/shutter/shutter.go

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -race -cover .

# Set GITHUB_TOKEN and create release git tag
.PHONY: release
release:
	goreleaser --rm-dist
