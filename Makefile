#                                                                                                                                                                                                                                                                                                                                                                                                                                                            #  SYSADM Server
#  Author  Wayne Wang <net_use@bzhy.com>
#  Copyright Bzhy Network
#  @HomePage http://www.sysadm.cn
#  @Version 0.21.03
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#  http://www.apache.org/licenses/LICENSE-2.0
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#  @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html
#  @Modified Mar 26 2021
#    

GO = GO111MODULE=on go
SERVER_GO_FILES ?= $(shell find ./src/server ! -name *_test.go  -name *.go -type f )
TEST_LOGGER_GO_FILES ?= $(shell find src/logger -name *_test.go  -type f )
SH_FILES ?= $(shell find ./scripts -name *.sh)
GOVERSION ?= $(shell go version)
BUILDTIME ?= $(shell date +'%Y.%m.%d.%H%M%S')
GITCOMMIT ?= $(shell git log --pretty=oneline -n 1)
BranchInfo ?= $(shell git rev-parse --abbrev-ref HEAD)

LDFlags=" \
    -X 'config.Commit=${GITCOMMIT}' \
    -X 'config.BuildBranch=${BranchInfo}' \
    -X 'config.Buildstamp=${BUILDTIME}' \
    -X 'config.goversion=${GOVERSION}' \
"
include Makefile.common

.PHONY: all server install clean test_logger

all: server

server:  ##Build  server
	@echo "building server"
	$(GO) build -ldflags $(LDFlags)  -o ./bin/server  $(SERVER_GO_FILES) 

install: ## Installing files to destination path
	@echo "Installing files to destination path"
	@chmod +x ./build/install.sh
	@./build/install.sh

clean:  ## Clean up intermediate build artifacts.
	@echo "cleaning" 
	@rm -rf ./bin/*
	@echo "$(INSTALL_DEST_PATH)"
#	@rm -rf /usr/local/

