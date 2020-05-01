BINARY_NAME := knihomolapp
DOCKER_REPOSITORY := danielkraic
DOCKER_IMAGE := knihomol

VERSION := $(shell git describe --tags 2> /dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD)
BUILD := $(shell date +"%FT%T%z")

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Commit=$(COMMIT) -X=main.Build=$(BUILD)"

all: test race cover build

build: fmt lint vet errcheck build-only

build-only:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME) -v

clean:
	go clean
	rm -f $(BINARY_NAME)

test: fmt vet lint errcheck
	go test -v ./...

race:
	go test -race ./...

cover:
	go test -cover ./...

coverprofile:
	go test -coverprofile=coverage.out  ./...
	go tool cover -html=coverage.out

errcheck:
	go list ./... | xargs errcheck

vet:
	go list ./... | grep -v vendor | xargs go vet

lint:
	go list ./... | xargs golint

fmt:
	go fmt ./...

docker:
	docker build -t $(DOCKER_REPOSITORY)/$(DOCKER_IMAGE):$(VERSION) .
