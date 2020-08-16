GO    	 := go
pkgs      = $(shell $(GO) list ./... | grep -v /vendor/)
arch      = amd64  ## default architecture
platforms = darwin linux
package   = batch_proc

PREFIX                  ?= $(shell pwd)
BIN_DIR                 ?= $(shell pwd)
DOCKER_REPO             ?= rolandvarga
DOCKER_IMAGE_NAME       ?= $(package)
DOCKER_IMAGE_TAG        ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

all: vet format test build generate

build: ## build executable for current platform
	@echo ">> building..."
	@$(GO) build

xbuild: ## cross build executables for all defined platforms
	@echo ">> cross building executable(s)..."

	@for platform in $(platforms); do \
		echo "build for $$platform/$(arch)" ;\
		name=$(package)'-'$$platform'-'$(arch) ;\
		echo $$name ;\
		GOOS=$$platform GOARCH=$(arch) $(GO) build -o $$name . ;\
	done

docker: all ## build docker image
	@echo ">> building docker image"
	@docker build -t "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

test: ## test code
	@echo ">> running tests.."
	@$(GO) test -v $(pkgs)

format: ## format code
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet: ## vet code
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

generate: ## generate test data
	@chmod +x $(PREFIX)/generate_data.sh
	$(PREFIX)/generate_data.sh

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
