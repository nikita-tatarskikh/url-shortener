# COMMON VARS
CONFIG_PROTOBUF_PATH=./cfg

MAKEFLAGS+=--silent

OS=$(shell uname -o)
PROJECTNAME=$(shell basename "$(PWD)")
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILDTIME=$(shell date -u +%Y-%m-%d_%H:%M:%S)
GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

ifeq ($(OS), Msys)
GOOS=windows
endif

# COLORED OUTPUT

GREEN=
LGREEN=
YELLOW=
ORANGE=
NC=# No Color

ifeq (${OS}, GNU/Linux)
GREEN=\033[1;32m
LGREEN=\033[0;32m
YELLOW=\033[1;33m
ORANGE=\033[0;33m
NC=\033[0m # No Color
endif

.PHONY: help
all: help
help:
	@echo
	@echo " Choose a command to run in "$(PROJECTNAME)":"
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | column -t -s ':' | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo

.PHONY: fmt generate lint go-mod-verify pre-commit benchmark build build-for-docker swagger-gen protoc-gen test test-with-integration coverage compose-test-up compose-test-down compose-test-down-with-volumes clean

## go-mod-verify: clean and verify go modules
go-mod-verify:
	@echo ">${YELLOW} Fixing modules...${NC}"
	@echo "  >${ORANGE} Adding missing and removing unused modules...${NC}"
	go mod tidy
	@echo "  >${ORANGE} Verifying dependencies have expected content...${NC}"
	go mod verify
	@echo ">${GREEN} Modules fixed${NC}"

## fmt: format all source files
fmt:
	@echo ">${YELLOW} Formating source files...${NC}"
	go fmt ./...
	@echo ">${GREEN} Source files formatted${NC}"

## generate: run all go generate commands from source files
generate:
	@echo ">${YELLOW} Generating source files...${NC}"
	go generate ./...
	@echo ">${GREEN} Source files generated${NC}"

## lint: perform static code analysis with golangci-lint tool (more than 30 linters inside)
lint:
	@echo ">${YELLOW} Linting source files...${NC}"
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix
	@echo ">${GREEN} Source files fine${NC}"

## benchmark: run all benchmarks inside project
benchmark:
	go test -bench=. ./...

## pre-commit: make sure the commit is safe
pre-commit: go-mod-verify fmt generate build test-with-integration lint
	@echo ">${GREEN} Commit can be made${NC}"

## swagger-gen: generate swagger specification for adserver-api application
swagger-gen:
	@echo ">${YELLOW} Generating swagger specification from source files...${NC}"
	go run github.com/swaggo/swag/cmd/swag init --parseDependency --parseInternal
	@echo ">${GREEN} Swagger specification is generated${NC}"

# build variables
BUILD_DIR=.
BUILDVARS=GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED}
DOCKER_BUILDVARS=GOOS=linux GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED}
XFLAGS=
LDFLAGS=-ldflags "-s -w ${XFLAGS}"
BUILD_CMD=${BUILDVARS} go build ${LDFLAGS}
DOCKER_BUILD_CMD=${DOCKER_BUILDVARS} go build ${LDFLAGS}

## build-adserver-api: build adserver-api application
build:
	@echo ">${YELLOW} Building url-shortener...${NC}"
	${BUILD_CMD} -o ${BUILD_DIR}/url_shortener ./cmd
	@echo ">${GREEN} url-shortener is built${NC}"

build-for-docker:
	@echo ">${YELLOW} Building url-shortener for docker...${NC}"
	${DOCKER_BUILD_CMD} -o ${BUILD_DIR}/url_shortener ./cmd
	@echo ">${GREEN} url-shortener is built for docker${NC}"

## test: run unit tests with race detection
test:
	@echo ">${YELLOW} Running unit tests with race detection...${NC}"
	go test -race -v ./...
	@echo ">${GREEN} All tests passed${NC}"

## test-with-integration: run unit and integration tests with race detection, docker required.
test-with-integration:
	@echo ">${YELLOW} Running unit and integration tests with race detection...${NC}"
	DOCKER_TEST=true go test -race -v ./...
	@echo ">${GREEN} All tests passed${NC}"


# docker compose variables
COMPOSE_TEST_FILE=docker-compose.yml
COMPOSE_TEST_CMD=docker-compose --project-name ${PROJECTNAME} --file ${COMPOSE_TEST_FILE}
COMPOSE_TEST_PULL_CMD=${COMPOSE_TEST_CMD} pull

## compose-test-up: raise the whole project from docker-compose.yml
compose-test-up: build-for-docker
	@echo ">${YELLOW} Raise the whole project from docker-compose.yml...${NC}"
	${COMPOSE_TEST_PULL_CMD}
	${COMPOSE_TEST_CMD} up --build --detach
	@echo ">${GREEN} Project raised${NC}"

## compose-test-down: destroy everything raised from docker-compose.yml
compose-test-down:
	@echo ">${YELLOW} Destroying everything raised from docker-compose.yml...${NC}"
	${COMPOSE_TEST_CMD} down --remove-orphans
	@echo ">${GREEN} Everything destroyed${NC}"

## compose-test-down-with-volumes: destroy everything raised from docker-compose.yml with volumes
compose-test-down-with-volumes:
	@echo ">${YELLOW} Destroying everything raised from docker-compose.yml with volumes...${NC}"
	${COMPOSE_TEST_CMD} down --remove-orphans --volumes
	@echo ">${GREEN} Everything destroyed${NC}"

## clean: run cache, modcache and testcache cleaning
clean: clean-cache clean-modcache clean-testcache

## clean-cache: remove the entire go build cache
clean-cache:
	@echo ">${YELLOW} Removing the entire go build cache...${NC}"
	go clean --cache ./...
	@echo ">${GREEN} Build cache removed${NC}"

## clean-modcache: remove the entire module download cache, including unpacked source code of versioned dependencies
clean-modcache:
	@echo ">${YELLOW} Removing the entire module download cache, including unpacked source code of versioned dependencies...${NC}"
	go clean --modcache ./...
	@echo ">${GREEN} Module download cache removed${NC}"

## clean-testcache: expire all test results in the go build cache
clean-testcache:
	@echo ">${YELLOW} Expiring all test results in the go build cache...${NC}"
	go clean --testcache ./...
	@echo ">${GREEN} All test results is expired${NC}"

## clean-docker: run Docker garbage collection of containers, images and volumes
clean-docker:
	@echo ">${YELLOW} Running docker garbage collection...${NC}"
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v /etc:/etc -e REMOVE_VOLUMES=1 spotify/docker-gc
	@echo ">${GREEN} Docker garbage successfully collected${NC}"