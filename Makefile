.PHONY: dev setup build install image test release clean

CGO_ENABLED=0
VERSION=$(shell git describe --abbrev=0 --tags)
COMMIT=$(shell git rev-parse --short HEAD)

all: dev

dev: build
	@./shops -v

build: clean
	@go build \
		-tags "netgo static_build" -installsuffix netgo 	\
		-ldflags "-w -X main.Version=$(VERSION) \
					 -X main.Commit=$(COMMIT)"	\
		.

install: build
	@go install

image:
	@docker build -t prologic/shops .

test: install
	@go test

release:
	@./tools/release.sh

clean:
	@git clean -f -d -X
