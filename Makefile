SOURCE_DIR = .
SOURCES := $(shell find $(SOURCE_DIR) -name *.go)

BINARY=hcsg

VERSION = 0.0.1
BUILD_TIME = $(shell date -u '+%Y-%m-%d %I:%M:%S %Z')
GIT_HASH = $(shell git rev-parse HEAD)
DATA_FILES := $(shell find conf | sed 's/ /\\ /g')
RELEASE_ROOT := "release"
TAGS = ""

LDFLAGS += -X "github.com/credli/hcsg/settings.AppVer=${VERSION}"
LDFLAGS += -X "github.com/credli/hcsg/settings.BuildTime=${BUILD_TIME}"
LDFLAGS += -X "github.com/credli/hcsg/settings.BuildGitHash=${GIT_HASH}"

.DEFAULT_GOAL: $(BINARY)

.PHONY: bindata install clean

$(BINARY): $(SOURCES)
	@rm -rf bin/
	@mkdir -p bin/public
	@cp -r "public/" bin/public/
	@go-bindata -pkg bindata -o bindata/bindata.go conf/
	@go build -ldflags '$(LDFLAGS)' -tags '$(TAGS)' -o bin/${BINARY}
	@echo Build complete

bindata: bindata/bindata.go

bindata/bindata.go: $(DATA_FILES)
	go-bindata -o=$@ -ignore="\\.DS_Store|README.md" -pkg=bindata conf/...

less: public/css/hcsg.css

public/css/hcsg.css: $(LESS_FILES)
	lessc $< $@

install:
	go install -ldflags ${LDFLAGS} ./...

clean:
	if [ -f bin/${BINARY} ] ; then rm bin/${BINARY} ; fi
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -print0 | xargs -0 rm