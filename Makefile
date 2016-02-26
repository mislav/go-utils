GOCMD = go
NAMESPACE = github.com/mislav/go-utils

TMPDIR ?= /tmp
GOPATH = $(TMPDIR)/$(subst /,_,$(NAMESPACE))
TMPLINK = $(GOPATH)/src/$(NAMESPACE)
PACKAGES = $(shell find . -depth 2 -name '*.go' | cut -d/ -f2 | sort -u)

.PHONY: test fmt

all: test

$(TMPLINK):
	mkdir -p $(dir $@)
	ln -snf $$PWD $@

test: $(TMPLINK)
	GO15VENDOREXPERIMENT=1 GOPATH=$(GOPATH) $(GOCMD) test $(patsubst %,$(NAMESPACE)/%,$(PACKAGES))

fmt:
	$(GOCMD) fmt $(patsubst %,./%,$(PACKAGES))
