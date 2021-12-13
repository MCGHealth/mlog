default: test

GO_ENUM_VERSION=v0.3.9
OS=$(shell uname -s)
ARCH = $(shell uname -m)

build:
	go build ./...

test:
	if [ ! -d "_coverage/" ]; then mkdir "_coverage/";fi
	go test ./... -coverprofile=_coverage/c1.tmp
	cat _coverage/c1.tmp | grep -v "logeventlevel_enum.go" > _coverage/coverage.out
	go tool cover -html=_coverage/coverage.out

local_dir:
	if [ ! -d "_local/" ]; then mkdir "_local/";fi
	curl -fsSL "https://github.com/abice/go-enum/releases/download/$(GO_ENUM_VERSION)/go-enum_$(OS)_$(ARCH)" -o _local/go-enum
	chmod 700 _local/go-enum

enum: local_dir	
	go generate ./...
	