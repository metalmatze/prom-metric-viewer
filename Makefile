EXECUTABLE ?= prom-metric-viewer
IMAGE ?= metalmatze/$(EXECUTABLE)
GO := CGO_ENABLED=0 go
DATE := $(shell date -u '+%FT%T%z')

LDFLAGS += -X main.Version=$(DRONE_TAG)
LDFLAGS += -X main.Revision=$(DRONE_COMMIT)
LDFLAGS += -X "main.Date=$(DATE)"
LDFLAGS += -extldflags '-static'

PACKAGES = $(shell go list ./... | grep -v /vendor/)

.PHONY: all
all: build

.PHONY: deps
deps:
	cd ui && pub get

.PHONY: clean
clean:
	$(GO) clean -i ./...
	packr clean
	rm -rf public/build.js

.PHONY: ui
ui:
	cd ui && pub build

.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/golang/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: test
test:
	@for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

.PHONY: packr
packr:
	@which packr > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/gobuffalo/packr/packr; \
	fi
	packr

$(EXECUTABLE): $(wildcard *.go)
	$(GO) build -v -ldflags '-w $(LDFLAGS)'

.PHONY: build
build: packr $(EXECUTABLE)

.PHONY: install
install:
	$(GO) install -v -ldflags '-w $(LDFLAGS)'

.PHONY: release
release: packr
	@which gox > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mitchellh/gox; \
	fi
	CGO_ENABLED=0 gox -verbose -ldflags '-w $(LDFLAGS)' -output="dist/$(EXECUTABLE)-${DRONE_TAG}-{{.OS}}-{{.Arch}}"
