.PHONY: all
all: help

# set default as dev if not set
export commit ?= HEAD

.PHONY: build

install: ## install project dependences
	go get -v ./...

tests: ## run go tests
	go test -v -race ./...

coverage: ## outputs coverage to coverage.out
	which goverage &> /dev/null || go get -u github.com/haya14busa/goverage
	goverage -v -race -coverprofile=coverage.out ./...

send-statiscs: ## send statistics to code quality services
	bash -c "$$(curl -s https://codecov.io/bash)"
	which goverage &> /dev/null || go get -u github.com/schrej/godacov
	godacov -t ${CODACY_TOKEN} -r ./coverage.out -c $(commit)

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'