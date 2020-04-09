PREFIX            ?= $(shell pwd)
FILES_TO_FMT      ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)

# Ensure everything works even if GOPATH is not set, which is often the case.
# The `go env GOPATH` will work for all cases for Go 1.8+.
GOPATH            ?= $(shell go env GOPATH)
TMP_GOPATH        ?= /tmp/sma-solar-exporter-go
GOBIN             ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO111MODULE       ?= on
export GO111MODULE

GOLANGCILINT_VERSION ?= d2b1eea2c6171a1a1141a448a745335ce2e928a1
GOLANGCILINT         ?= $(GOBIN)/golangci-lint-$(GOLANGCILINT_VERSION)
MISSPELL_VERSION     ?= c0b55c8239520f6b5aa15a0207ca8b28027ba49e
MISSPELL             ?= $(GOBIN)/misspell-$(MISSPELL_VERSION)
FAILLINT_VERSION     ?= v1.2.0
FAILLINT             ?= $(GOBIN)/faillint-$(FAILLINT_VERSION)
YASDI_URL            ?= https://files.sma.de/dl/11705/yasdi-1.8.1build9-src.zip
YASDI                ?= $(PREFIX)/yasdi

# fetch_go_bin_version downloads (go gets) the binary from specific version and installs it in $(GOBIN)/<bin>-<version>
# arguments:
# $(1): Install path. (e.g github.com/campoy/embedmd)
# $(2): Tag or revision for checkout.
define fetch_go_bin_version
	@mkdir -p $(GOBIN)
	@mkdir -p $(TMP_GOPATH)

	@echo ">> fetching $(1)@$(2) revision/version"
	@if [ ! -d '$(TMP_GOPATH)/src/$(1)' ]; then \
    GOPATH='$(TMP_GOPATH)' GO111MODULE='off' go get -d -u '$(1)/...'; \
  else \
    CDPATH='' cd -- '$(TMP_GOPATH)/src/$(1)' && git fetch; \
  fi
	@CDPATH='' cd -- '$(TMP_GOPATH)/src/$(1)' && git checkout -f -q '$(2)'
	@echo ">> installing $(1)@$(2)"
	@GOBIN='$(TMP_GOPATH)/bin' GOPATH='$(TMP_GOPATH)' GO111MODULE='off' go install '$(1)'
	@mv -- '$(TMP_GOPATH)/bin/$(shell basename $(1))' '$(GOBIN)/$(shell basename $(1))-$(2)'
	@echo ">> produced $(GOBIN)/$(shell basename $(1))-$(2)"

endef

.PHONY: all
all: deps lint build

.PHONY: build
build: ## Builds binaries
build: deps
	@echo ">> building binaries $(GOBIN)"
	go build ./cmd/sma-solar-exporter
	go build ./cmd/exporter-simulation

.PHONY: deps
deps: yasdi
	## Ensures fresh go.mod and go.sum.
	@go mod tidy
	@go mod verify

lint: $(FAILLINT) $(GOLANGCILINT) $(MISSPELL)
	@echo ">> verifying modules being imported"
		@$(FAILLINT) -paths "errors=github.com/pkg/errors,\
		github.com/prometheus/tsdb=github.com/prometheus/prometheus/tsdb,\
		github.com/prometheus/prometheus/pkg/testutils=github.com/thanos-io/thanos/pkg/testutil,\
		github.com/prometheus/client_golang/prometheus.{DefaultGatherer,DefBuckets,NewUntypedFunc,UntypedFunc},\
		github.com/prometheus/client_golang/prometheus.{NewCounter,NewCounterVec,NewCounterVec,NewGauge,NewGaugeVec,NewGaugeFunc,\
		NewHistorgram,NewHistogramVec,NewSummary,NewSummaryVec}=github.com/prometheus/client_golang/prometheus/promauto.{NewCounter,\
		NewCounterVec,NewCounterVec,NewGauge,NewGaugeVec,NewGaugeFunc,NewHistorgram,NewHistogramVec,NewSummary,NewSummaryVec}" ./...
	@echo ">> linting all of the Go files GOGC=${GOGC}"
	@$(GOLANGCILINT) run
	@echo ">> detecting misspells"
	@find . -type f | grep -v yasdi | grep -vE '\./\..*' | xargs $(MISSPELL) -error
	@echo ">> examining all of the Go files"
	@go vet -stdmethods=false ./pkg/... ./cmd/...

yasdi:
	@mkdir -p yasdi
	
	@echo ">> fetching yasdi from $(YASDI_URL)"
	@curl --output yasdi/yasdi.zip $(YASDI_URL)
	@echo ">> unpacking yasdi/yasdi.zip"
	@CDPATH='' cd -- 'yasdi' && unzip yasdi.zip && rm -rf yasdi.zip
	
	@mkdir -p yasdi/projects/generic-cmake/build-gcc
	@CDPATH='' cd -- 'yasdi/projects/generic-cmake/build-gcc' && cmake .. && make
	
	@echo ">> yasdi is now prepared"

$(GOLANGCILINT):
	$(call fetch_go_bin_version,github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCILINT_VERSION))

$(MISSPELL):
	$(call fetch_go_bin_version,github.com/client9/misspell/cmd/misspell,$(MISSPELL_VERSION))

$(FAILLINT):
	$(call fetch_go_bin_version,github.com/fatih/faillint,$(FAILLINT_VERSION))
