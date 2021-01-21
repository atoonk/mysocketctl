GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=mysocketctl
BUCKET=mysocketctl

all: test build

release:
	GOOS=windows GOARCH=amd64 go build -o ./bin/$(BINARY_NAME)_windows_amd64
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(BINARY_NAME)_linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o ./bin/$(BINARY_NAME)_darwin_amd64
	python3 ./s3upload.py ./bin/mysocketctl_darwin_amd64 ${BUCKET} macos/mysocketctl
	python3 ./s3upload.py ./bin/mysocketctl_linux_amd64 ${BUCKET} linux_amd64/mysocketctl
	python3 ./s3upload.py ./bin/mysocketctl_windows_amd64 ${BUCKET} windows_amd64/mysocketctl


build:
	$(GOBUILD) -o $(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

