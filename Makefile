PKG := github.com/patrislav/marwind

.PHONY: all
all: bin/marwm

install:


bin/marwm:
	go build -o bin/marwm $(PKG)
