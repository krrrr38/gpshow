BIN = gpshow

build: deps
	go build -o $(BIN)

lint: deps
	go vet
	golint

test: testdeps
	go test ./...

deps:
	go get -d -v github.com/jteeuwen/go-bindata
	go get -d -v .
	go-bindata -o=resources.go ./resources/...

testdeps:
	go get -d -v -t .

install: deps
	go install

clean:
	rm resources.go
	go clean
