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

CONFIG_DIR := $(DESTDIR)/etc/secpaver
SCRIPTS_DIR := $(DESTDIR)/usr/share/secpaver/scripts

BUILDFLAGS := -trimpath
LDFLAGS := -w -s -buildid=IdBySecPaver -linkmode=external -extldflags=-static -extldflags=-zrelro -extldflags=-Wl,-z,now

ifeq ($(shell go help mod >/dev/null 2>&1 && echo true), true)
export GO111MODULE=on
export GOFLAGS=-mod=vendor
endif

all: sec_conf

sec_conf:
	CGO_CFLAGS_ALLOW="-ftrapv -D_FORTIFY_SOURCE=2 -O2" CGO_CFLAGS="-fstack-protector-strong -ftrapv -D_FORTIFY_SOURCE=2 -O2" \
	CGO_LDFALGS_ALLOW="-Wl,-z,-s,relro,now,noexecstack" CGO_LDFALGS="-Wl,-z,-s,relro,now,noexecstack" \
	go build -buildmode=pie $(BUILDFLAGS) -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/sec_conf ./sec_conf.go

test:
	CGO_CFLAGS_ALLOW="-ftrapv -D_FORTIFY_SOURCE=2 -O2" CGO_CFLAGS="-fstack-protector-strong -ftrapv -D_FORTIFY_SOURCE=2 -O2" \
	CGO_LDFALGS_ALLOW="-Wl,-z,-s,relro,now,noexecstack" CGO_LDFALGS="-Wl,-z,-s,relro,now,noexecstack" \
	go test -buildmode=pie $(BUILDFLAGS) -ldflags '$(LDFLAGS)' secconf/*.go

everything: sec_conf

clean:
	rm -f $(BUILD_DIR)/sec_conf

install:
	@echo "BEGIN INSTALL sec_conf"

	mkdir -p $(BIN_DIR)
	mkdir -p $(SCRIPTS_DIR)
	mkdir -p $(SCRIPTS_DIR)/sec_conf
	mkdir -p $(SCRIPTS_DIR)/sec_conf/check
	mkdir -p $(SCRIPTS_DIR)/sec_conf/gen

	install -m 400 secconf/gen_comm.sh $(SCRIPTS_DIR)/sec_conf
	install -m 600 secconf/sec_conf.yaml $(SCRIPTS_DIR)/sec_conf
	install -m 400 secconf/gen/* $(SCRIPTS_DIR)/sec_conf/gen
	install -m 400 secconf/check/* $(SCRIPTS_DIR)/sec_conf/check
	install -m 500 $(BUILD_DIR)/sec_conf $(BIN_DIR)

	@echo "END INSTALL sec_conf"

