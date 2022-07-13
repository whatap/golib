
test:
	go test ./... -v

build:
	go mod download -x
	go build ./...

clean :
	go clean 
