SRC=*.go cmd/terminal-to-html/*.go
BINARY=terminal-to-html
BUILDCMD=go build -o $@ cmd/terminal-to-html/*
VERSION=$(shell cat version.go  | grep baseVersion | head -n1 | cut -d \" -f 2)

all: test $(BINARY)

bench:
	go test -bench . -benchmem

test:
	go test

clean:
	rm -f $(BINARY)
	rm -rf dist bin

cmd/terminal-to-html/_bindata.go: assets/terminal.css
	go-bindata -o cmd/terminal-to-html/bindata.go -nomemcopy assets

$(BINARY): $(SRC)
	$(BUILDCMD)

version:
	@echo $(VERSION)

# Cross-compiling

GZ_ARCH     := linux-amd64 linux-i386 linux-armel darwin-i386 darwin-amd64
ZIP_ARCH    := windows-i386 windows-amd64
GZ_TARGETS  := $(foreach target,$(GZ_ARCH), dist/$(BINARY)-$(VERSION)-$(target).gz)
ZIP_TARGETS := $(foreach target,$(ZIP_ARCH), dist/$(BINARY)-$(VERSION)-$(target).zip)

dist: $(GZ_TARGETS) $(ZIP_TARGETS)

dist/%.gz: bin/%
	@[ -d dist ] || mkdir dist
	gzip -c $< > $@

dist/%.zip: bin/%
	@[ -d dist ] || mkdir dist
	@rm -f $@ || true
	zip $@ $<

bin/$(BINARY)-$(VERSION)-%: $(SRC)
	@[ -d bin ] || mkdir bin
	GOOS=$(firstword $(subst -, , $*)) GOARCH=$(lastword $(subst armel, arm, $(subst i386, 386, $(subst -, , $*)))) $(BUILDCMD)

.PHONY: clean bench test dist version
