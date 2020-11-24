ROOT ?= $(shell pwd)

.PHONY: all

all: build

.PHONY: protos
protos:
	@sh script/protos.sh

.PHONY: build
build:
	mkdir -p builds/debug
	go build -o builds/debug/service -ldflags '${EXTLDFLAGS}-X dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/version.VersionDev=build.$(BUILD_NUMBER)' dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher/example

.PHONY: run
run:build
	builds/debug/service --config=./config/example.yaml
