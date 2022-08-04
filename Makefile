
build: package test

generate:
	go generate ./...

test:
	go test ./... -cover

package:
	go mod download -x
	go build ./...

clean :
	go clean 