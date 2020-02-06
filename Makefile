GOVERSION:=$(shell go version)
GOOS:=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH:=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
BUILD_DIR:=build/$(GOOS)-$(GOARCH)
REVISION:= $(shell git rev-parse HEAD)

APPNAME:=sideupload
REVISIONPACK:=github.com/neilli-sable/sideupload.revision=$(REVISION)

.PHONY: all generate install update build compress serve clean clean-with-vendor

all: build compress docker

generate:
	go generate ./...

install:

update:
	go get -u

build:
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)
	go build -ldflags="-s" -o $(BUILD_DIR)/bin/$(APPNAME)

compress:
	upx $(BUILD_DIR)/bin/$(APPNAME)

serve:
	cd dependency/ && docker-compose up -d

serve-down:
	cd dependency/ && docker-compose down

cron:
	$(BUILD_DIR)/bin/$(APPNAME) cron

docker: build compress
	docker build -t $(APPNAME):latest ./ 

test:
	go test -race -coverprofile=profile.out  ./...
	go tool cover -html=profile.out  -o cover.html

clean:
	rm -rf build package

clean-with-vendor:
	rm -rf build package vendor .vendor-new
