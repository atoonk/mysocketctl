GOCMD=go
GOBUILD=$(GOCMD) build
GOFMT=$(GOCMD)fmt
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=mysocketctl
BUCKET=mysocketctl

DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(shell git describe --long --dirty --tags)
FLAGS := -ldflags "-X github.com/atoonk/mysocketctl/go/cmd.version=$(VERSION) -X github.com/atoonk/mysocketctl/go/cmd.date=$(DATE)"

all: lint test build

release:
	GOOS=windows GOARCH=amd64 go build $(FLAGS) -o ./bin/$(BINARY_NAME)_windows_amd64
	GOOS=linux GOARCH=amd64 go build $(FLAGS) -o ./bin/$(BINARY_NAME)_linux_amd64
	GOOS=linux GOARCH=arm64 go build $(FLAGS) -o ./bin/$(BINARY_NAME)_linux_arm64
	GOOS=darwin GOARCH=amd64 go build $(FLAGS)  -o ./bin/$(BINARY_NAME)_darwin_amd64

	shasum -a 256 ./bin/mysocketctl_darwin_amd64 | awk '{print $$1}' > ./bin/mysocketctl_darwin_amd64-sha256-checksum.txt
	python3 ./s3upload.py ./bin/mysocketctl_darwin_amd64-sha256-checksum.txt ${BUCKET} darwin_amd64/sha256-checksum.txt
	python3 ./s3upload.py ./bin/mysocketctl_darwin_amd64 ${BUCKET} darwin_amd64/mysocketctl

	shasum -a 256 ./bin/mysocketctl_linux_amd64 | awk '{print $$1}' > ./bin/mysocketctl_linux_amd64-sha256-checksum.txt
	python3 ./s3upload.py ./bin/mysocketctl_linux_amd64-sha256-checksum.txt ${BUCKET} linux_amd64/sha256-checksum.txt
	python3 ./s3upload.py ./bin/mysocketctl_linux_amd64 ${BUCKET} linux_amd64/mysocketctl

	#This is for Raspberrypi
	shasum -a 256 ./bin/mysocketctl_linux_arm64 | awk '{print $$1}' > ./bin/mysocketctl_linux_arm64-sha256-checksum.txt
	python3 ./s3upload.py ./bin/mysocketctl_linux_arm64-sha256-checksum.txt ${BUCKET} linux_arm64/sha256-checksum.txt
	python3 ./s3upload.py ./bin/mysocketctl_linux_arm64 ${BUCKET} linux_arm64/mysocketctl

	shasum -a 256 ./bin/mysocketctl_windows_amd64 | awk '{print $$1}' > ./bin/mysocketctl_windows_amd64-sha256-checksum.txt
	python3 ./s3upload.py ./bin/mysocketctl_windows_amd64-sha256-checksum.txt ${BUCKET} windows_amd64/sha256-checksum.txt
	python3 ./s3upload.py ./bin/mysocketctl_windows_amd64 ${BUCKET} windows_amd64/mysocketctl.exe

	echo ${VERSION} > latest_version.txt
	python3 ./s3upload.py latest_version.txt ${BUCKET} latest_version.txt
	rm latest_version.txt


build:
	$(GOBUILD) $(FLAGS) -o $(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(FLAGS) -o $(BINARY_NAME) -v

lint:
	@echo "running go fmt"
	$(GOFMT) -w .

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) $(FLAGS) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

