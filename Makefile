SHELL := /bin/bash
MODULE = $(shell go list -m)
PACKAGES := $(shell go list ./... | grep -v /vendor/)

.PHONY: default
default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' 


.PHONY: build
build:  ## build the spellcorrect API server binary
	CGO_ENABLED=0 go build -a -o spell-correct-server $(MODULE)

.PHONY: clean
clean: ## remove temporary files
	rm -rf spell-correct-server coverage.out coverage-all.out

.PHONY: test
test: ## run unit tests
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)


.PHONY: docker-build
docker-build:  ## build the spellcorrect API server docker image
	docker build -t spell-correctort .


.PHONY: docker-run
docker-run:  ## runs the spellcorrect API server docker image
	docker run  -p 10000:10000 -v /home/giorgos/Development/github.com/gosom/context-spell-correct/datasets:/datasets -e SC_SENTENCES_PATH=/datasets/sentences.txt2 -e SC_DICT_PATH=/datasets/de-100k.txt spell-correctort 

