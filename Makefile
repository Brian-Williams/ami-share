### AMI share utility

MKFILE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export GOPATH=$(MKFILE_PATH)/..

platform=$(shell uname -s)
ifeq ($(platform), Darwin)
	GOOS=Darwin
	GO = $(shell command -v go)
else ifeq ($(platform), Linux)
	GO = /usr/local/bin/go/bin/go
else
	GOOS = Windows
	GO = "c:\Program Files (x86)\Go\bin\go"
endif

build: install

install:
	$(GO) install ./
