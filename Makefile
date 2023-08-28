# Copyright (c) Huawei Technologies Co., Ltd. 2020-2021. All rights reserved.
# secPaver is licensed under the Mulan PSL v2.
# You can use this software according to the terms and conditions of the Mulan PSL v2.
# You may obtain a copy of Mulan PSL v2 at:
#     http://license.coscl.org.cn/MulanPSL2
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
# PURPOSE.
# See the Mulan PSL v2 for more details.

.PHONY: all clean install unit-test

PWD := $(shell pwd)
BUILD_DIR := $(PWD)/build
BIN_DIR := $(DESTDIR)/usr/bin
SO_DIR := $(DESTDIR)/usr/lib64/secpaver

CONFIG_DIR := $(DESTDIR)/etc/secpaver
RESOURCE_DIR := $(DESTDIR)/var/local/secpaver
SCRIPTS_DIR := $(DESTDIR)/usr/share/secpaver/scripts
DATA_DIR := $(DESTDIR)/usr/share/secpaver
SYSTEMD_DIR = $(DESTDIR)/usr/lib/systemd/system
LOG_DIR := $(DESTDIR)/var/log/secpaver

BUILDFLAGS := -trimpath
GO_LDFLAGS := -w -s -buildid=IdBySecPaver -linkmode=external -extldflags=-static -extldflags=-zrelro -extldflags=-znow

ifeq ($(shell go help mod >/dev/null 2>&1 && echo true), true)
export GO111MODULE=on
export GOFLAGS=-mod=vendor
endif

all: pav pavd

selinux:
	CGO_CFLAGS_ALLOW="-ftrapv -D_FORTIFY_SOURCE=2 -O2" CGO_CFLAGS="-fstack-protector-strong -ftrapv -D_FORTIFY_SOURCE=2 -O2" \
	CGO_LDFALGS_ALLOW="-Wl,-z,-s,relro,now,noexecstack" CGO_LDFALGS="-Wl,-z,-s,relro,now,noexecstack" \
	go build -buildmode=plugin $(BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' -o $(BUILD_DIR)/selinux.so cmd/plugin/selinux/*.go
	strip $(BUILD_DIR)/selinux.so

pav:
	CGO_CFLAGS_ALLOW="-ftrapv -D_FORTIFY_SOURCE=2 -O2" CGO_CFLAGS="-fstack-protector-strong -ftrapv -D_FORTIFY_SOURCE=2 -O2" \
	CGO_LDFALGS_ALLOW="-Wl,-z,-s,relro,now,noexecstack" CGO_LDFALGS="-Wl,-z,-s,relro,now,noexecstack" \
	go build -buildmode=pie $(BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' -o $(BUILD_DIR)/pav cmd/pav/*.go
	strip $(BUILD_DIR)/pav

pavd:
	CGO_CFLAGS_ALLOW="-ftrapv -D_FORTIFY_SOURCE=2 -O2" CGO_CFLAGS="-fstack-protector-strong -ftrapv -D_FORTIFY_SOURCE=2 -O2" \
	CGO_LDFALGS_ALLOW="-Wl,-z,-s,relro,now,noexecstack" CGO_LDFALGS="-Wl,-z,-s,relro,now,noexecstack" \
	go build -buildmode=pie $(BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' -o $(BUILD_DIR)/pavd cmd/pavd/*.go
	strip $(BUILD_DIR)/pavd

everything: pav pavd selinux

clean:
	rm -rf $(BUILD_DIR)/pav*

unit-test:
	go test -race -coverpkg ./... -coverprofile=coverage.data ./...
	go tool cover -func=coverage.data -o coverage.txt
	go tool cover -html=coverage.data -o coverage.html
	tail -1 coverage.txt

func-test:
	cd tests && chmod 700 *.sh && ./test.sh

install:
	@echo "BEGIN INSTALL secPaver"

	mkdir -p $(BIN_DIR)
	mkdir -p $(SO_DIR)
	mkdir -p $(SCRIPTS_DIR)
	mkdir -p $(SCRIPTS_DIR)/selinux
	mkdir -p $(SYSTEMD_DIR)
	mkdir -p $(LOG_DIR)
	mkdir -p $(CONFIG_DIR)/pavd
	mkdir -p $(RESOURCE_DIR)/projects
	mkdir -p $(RESOURCE_DIR)/policies/selinux

	chmod 700 $(RESOURCE_DIR) \
			  $(SO_DIR) \
			  $(LOG_DIR) \
			  $(CONFIG_DIR) $(CONFIG_DIR)/pavd \
			  $(RESOURCE_DIR)/projects $(RESOURCE_DIR)/policies $(RESOURCE_DIR)/policies/selinux \
			  $(SCRIPTS_DIR) $(SCRIPTS_DIR)/selinux \
			  $(DATA_DIR)

	install -m 500 $(BUILD_DIR)/pav $(BIN_DIR)
	install -m 500 $(BUILD_DIR)/pavd $(BIN_DIR)
	install -m 500 $(BUILD_DIR)/*.so $(SO_DIR)
	install -m 600 misc/config/config.json $(CONFIG_DIR)/pavd
	install -m 500 misc/scripts/selinux/*.sh $(SCRIPTS_DIR)/selinux
	install -m 600 misc/scripts/selinux/config $(SCRIPTS_DIR)/selinux
	install -m 600 misc/pavd.service $(SYSTEMD_DIR)
	systemctl daemon-reload

	@echo "END INSTALL secPaver"

