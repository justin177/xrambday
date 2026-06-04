APP := xrambday
ZIP := $(APP).zip
BOOTSTRAP := bootstrap
GOARCH ?= arm64
GOCACHE ?= /private/tmp/xrambday-gocache

.PHONY: build zip package package-amd64 package-arm64 clean

build:
	GOCACHE=$(GOCACHE) GOOS=linux GOARCH=$(GOARCH) CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BOOTSTRAP) .

zip: build
	zip -q -j $(ZIP) $(BOOTSTRAP)

package: zip

package-arm64:
	$(MAKE) package GOARCH=arm64

package-amd64:
	$(MAKE) package GOARCH=amd64

clean:
	rm -f $(BOOTSTRAP) $(ZIP)
