# Copyright (c) 2020 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(COMMIT)-$(DATE)

ifneq ($(shell git status --porcelain),)
    VERSION := $(VERSION)-dirty
endif

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -ldflags "-X main.version=$(VERSION)" -o ../build/plugin-unmanaged-racks/tmp/urp main.go
	cp ./config/config.yml ../build/plugin-unmanaged-racks/tmp/

clean:
	rm -r ../build/plugin-unmanaged-racks/tmp

run:
	PLUGIN_CONFIG_FILE_PATH=./config/config.yml go run ./main.go

lint: lint-deps
	golangci-lint run ./...

.PHONY: lint-deps
lint-deps:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.19.1
