GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=mysocketctl
BUCKET=mysocketctl

DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(shell git describe --long --dirty --tags)
FLAGS := -ldflags "-X github.com/atoonk/mysocketctl/go/cmd.version=$(VERSION) -X github.com/atoonk/mysocketctl/go/cmd.date=$(DATE)"

all: test build

release:
	GOOS=windows GOARCH=amd64 go build $(FLAGS) -o ./bin/$(BINARY_NAME)_windows_amd64
	GOOS=linux GOARCH=amd64 go build $(FLAGS) -o ./bin/$(BINARY_NAME)_linux_amd64
	GOOS=darwin GOARCH=amd64 go build $(FLAGS)  -o ./bin/$(BINARY_NAME)_darwin_amd64
	python3 ./s3upload.py ./bin/mysocketctl_darwin_amd64 ${BUCKET} macos/mysocketctl
	python3 ./s3upload.py ./bin/mysocketctl_linux_amd64 ${BUCKET} linux_amd64/mysocketctl
	python3 ./s3upload.py ./bin/mysocketctl_windows_amd64 ${BUCKET} windows_amd64/mysocketctl


build:
	$(GOBUILD) $(FLAGS) -o $(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(FLAGS) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) $(FLAGS) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

