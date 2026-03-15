VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS  = -s -w -X github.com/aymenhmaidiwastaken/devboot/cmd.Version=$(VERSION) \
           -X github.com/aymenhmaidiwastaken/devboot/cmd.Commit=$(COMMIT) \
           -X github.com/aymenhmaidiwastaken/devboot/cmd.Date=$(DATE)

.PHONY: build install clean test lint run

build:
	go build -ldflags "$(LDFLAGS)" -o devboot .

install:
	go install -ldflags "$(LDFLAGS)" .

clean:
	rm -f devboot devboot.exe

test:
	go test ./...

lint:
	golangci-lint run ./...

run:
	go run -ldflags "$(LDFLAGS)" . $(ARGS)

# Cross-compilation
.PHONY: build-all
build-all:
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/devboot-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/devboot-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/devboot-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/devboot-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/devboot-windows-amd64.exe .
