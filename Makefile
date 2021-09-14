EXTENSION:=
ifeq ($(OS),Windows_NT)
  EXTENSION:=.exe
endif

all: bin

LDFLAGS="-s -w -X $(PKG_NAME)/internal.Version=${GIT_TAG}"

bin:
	CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -o bin/docker-compose$(EXTENSION) ./main.go

cross:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o bin/docker-compose-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -ldflags=$(LDFLAGS) -o bin/docker-compose-linux-arm64 .

test:
	go test -cover $(shell go list  $(TAGS) ./... | grep -vE 'e2e')

test-ubuntu-install:
	docker build -f ubuntu-test.Dockerfile .

test-centos-install:
	docker build -f centos-test.Dockerfile .

test-install: test-centos-install test-ubuntu-install

.PHONY: bin cross test test-ubuntu-install test-centos-install
