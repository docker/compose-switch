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
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o bin/docker-compose-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -ldflags=$(LDFLAGS) -o bin/docker-compose-darwin-arm64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o bin/docker-compose-windows-amd64.exe .

test:
	go test -cover $(shell go list  $(TAGS) ./... | grep -vE 'e2e')

build-e2e-stub:
	go build -ldflags=$(LDFLAGS) -o bin/test/echostub$(EXTENSION) ./tests/echostub/main.go

e2e: build-e2e-stub
	go test -count=1 -v $(TEST_FLAGS) ./tests/e2e

.PHONY: bin cross test e2e
