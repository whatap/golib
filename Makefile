GO=go

build: package test

generate:
	$(GO) generate ./...

test:
	$(GO) test ./... -cover

package:
	$(GO) mod tidy
	$(GO) mod download -x
	$(GO) build ./...

clean :
	$(GO) clean -modcache
	$(GO) clean -testcache
	$(GO) clean -cache
	$(GO) clean
	rm -rf go.sum

version:
	$(GO) version
	$(GO) env
