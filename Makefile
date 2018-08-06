.PHONY: all
all: help

# set default as dev if not set
export commit ?= HEAD
export testWatchPort=8091

.PHONY: build

update-dev-deps: ## update dev tools
	go get -u -v github.com/haya14busa/goverage
	go get -u -v golang.org/x/lint/golint
	go get -u -v github.com/schrej/godacov

tests-watch:
	go get github.com/smartystreets/goconvey
	goconvey -port $(testWatchPort)

install: ## install project dependences
	go get -v ./...

lint: ## run got lint
	go get golang.org/x/lint/golint
	golint `find server -maxdepth 1 -type d`

tests: ## run go tests
	go test -v -race ./...

tests-watch: ## run got tests and keep watching for changes
	go get github.com/smartystreets/goconvey
	goconvey -port $(testWatchPort)

coverage: ## outputs coverage to coverage.out
	go get github.com/haya14busa/goverage
	goverage -v -race -coverprofile=coverage.out ./...

send-statiscs: ## send statistics to code quality services
	bash -c "$$(curl -s https://codecov.io/bash)"
	go get github.com/schrej/godacov
	godacov -t ${CODACY_TOKEN} -r ./coverage.out -c $(commit)

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

