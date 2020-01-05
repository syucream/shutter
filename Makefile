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

# Set GITHUB_TOKEN and create release git tag
.PHONY: release
release:
	goreleaser --rm-dist
