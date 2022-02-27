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

DATE := $(shell date +%Y%m%d)
VERSION := $(DATE)
BUILD_DIR := tmp
DEFAULT_ODIMRA_PASS := 'Od!m12$$4'

ENCRYPT_PASSWORD := echo -n ${DEFAULT_ODIMRA_PASS} \
		| openssl pkeyutl -encrypt -inkey ../cert_generator/odimra_rsa.private -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha512 \
		| openssl base64 -A

ifneq ($(shell git status --porcelain),)
    VERSION := $(VERSION)-dirty
endif

clean:
	rm -r BUILD_DIR

compile:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -ldflags "-X main.version=$(VERSION)" -o tmp/urp main.go

configure:
	@mkdir --parents ${BUILD_DIR}
	@sed "s|\#password\#|$(shell $(ENCRYPT_PASSWORD))|" config/config.yml > ${BUILD_DIR}/config.yml

build: configure compile

run: configure
	PLUGIN_CONFIG_FILE_PATH=${BUILD_DIR}/config.yml go run ./main.go

lint: lint-deps
	golangci-lint run ./...

.PHONY: lint-deps
lint-deps:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.19.1
