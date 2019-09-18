PKG := github.com/patrislav/marwind
SOURCES := $(shell find . -name '*.go')
VERSION := $(shell git describe --always --long --dirty)
BUILDTIME := $(shell date +'%Y-%m-%d_%T')
LDFLAGS :=

.PHONY: all
all: bin/marwm

bin/marwm: $(SOURCES)
	go build -o bin/marwm \
		-trimpath \
		-ldflags="-X main.version=$(VERSION) -X main.buildTime=$(BUILDTIME) $(LDFLAGS)" \
		$(PKG)/cmd/marwm
