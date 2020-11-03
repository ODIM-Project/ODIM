COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(COMMIT)-$(DATE)

ifneq ($(shell git status --porcelain),)
    VERSION := $(VERSION)-dirty
endif

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -ldflags "-X main.version=$(VERSION)" -o build/urp main.go

clean:
	rm -r ./build

run:
	PLUGIN_CONFIG_FILE_PATH=./config/config.json go run ./main.go

lint: lint-deps
	golangci-lint run ./...

.PHONY: lint-deps
lint-deps:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.19.1
