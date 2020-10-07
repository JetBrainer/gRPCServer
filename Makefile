BIN := main
.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...
## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-get

.PHONY: build-tokenizer
## build-tokenizer: build the tokenizer application
build-tokenizer:
	${MAKE} -c tokenizer build
.PHONY: build
## Build project: build command
build:
	cd cmd/apiserver && go build -o ../../bin/${BIN}
.PHONY: run
## Run project: run command
run:
	../../bin/${BIN}
.PHONY: clean
## Clean: clean command
clean:
	go clean
	clear
.PHONY: all
## Cleans binary file clears terminal and build/run project: all make file
all:
	cd cmd/apiserver && go build -o ../../bin/${BIN} && cd ../../ && ./bin/${BIN}

.PHONY: linter
## golangci-linter: go linter
lint:
	golangci-lint run
.PHONY: setup
## setup: setup go modules
setup:
	@go mod init \
		&& go mod tidy \
		&& go mod vendor

# helper rule for deployment
#check-environment:
#ifndef ENV
#    $(error ENV not set, allowed values - `staging` or `production`)
#endif

.PHONY: docker-build
## docker-build: builds the stringifier docker image to registry
docker-build: build
	docker build -t ${APP}:${COMMIT_SHA} .

.PHONY: docker-push
## docker-push: pushes the stringifier docker image to registry
docker-push: check-environment docker-build
	docker push ${REGISTRY}/${ENV}/${APP}:${COMMIT_SHA}

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'