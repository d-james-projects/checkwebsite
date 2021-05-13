VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAGS="-X main.Version=${VERSION} -X main.Build=${BUILD} -w -s"
TAGS?=v0.0
REPO?=dnjameshome

.PHONY: help
help: ## - Show help message
	@printf "\033[32m\xE2\x9c\x93 usage: make [target]\n\n\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build:	## - Build the golang docker image based on scratch
	@printf "\033[32m\xE2\x9c\x93 Build the golang docker image based on scratch\n\033[0m"
	@export DOCKER_CONTENT_TRUST=1 && docker build --build-arg flags=$(LDFLAGS) -f Dockerfile -t redgreenserver .

.PHONY: build-no-cache
build-no-cache:	## - Build the golang docker image based on scratch with no cache
	@printf "\033[32m\xE2\x9c\x93 Build the golang docker image based on scratch\n\033[0m"
	@export DOCKER_CONTENT_TRUST=1 && docker build --no-cache --build-arg flags=$(LDFLAGS) -f Dockerfile -t redgreenserver .

.PHONY: push-image
push-image:	## - Push docker image to container registry
	@docker tag redgreenserver:latest $(REPO)/redgreenserver:$(TAGS)
	@docker push $(REPO)/redgreenserver:$(TAGS)
