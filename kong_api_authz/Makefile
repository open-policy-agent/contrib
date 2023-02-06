BLUE=\033[1;36m
NC=\033[0m # No Color

.PHONY: all check build test verify clean push

# Execute the whole lifecycle
all: check build test verify

# Analyze code quality
check: docker-build-devcontainer
	@echo ""
	@echo "> ${BLUE}Checking Lua code...${NC}"
	@echo ""
	$(call docker-run-devcontainer,luacheck src spec)

# Build package in current directory using the rockspec
build: docker-build-devcontainer check
	@echo ""
	@echo "> ${BLUE}Installing package...${NC}"
	@echo ""
	$(call docker-run-devcontainer,luarocks make)

# Run all unit tests
test: docker-build-devcontainer check build
	@echo ""
	@echo "> ${BLUE}Running unit tests...${NC}"
	@echo ""
	$(call docker-run-devcontainer,luarocks test)

# Run integration tests
verify:
	@echo ""
	@echo "> ${BLUE}Running integration tests...${NC}"
	@echo ""
	$(call docker-compose-test,build)
	$(call docker-compose-test,run sut)
	$(call docker-compose-test,down -v --rmi local)

# Delete all files that are created by running make,
# including devcontainer image build for running Makefile commands
clean: mostlyclean docker-rm-devcontainer

# Like ‘clean’, but do not delete the Docker image built to run other commands
mostlyclean:
	@echo ""
	@echo "> ${BLUE}Cleaning up project directory...${NC}"
	@echo ""
	rm -rf .luarocks lua_modules lua luacov.* luarocks

# Publish a Kong image with the kong-plugin-opa plugin installed
push: docker-publish-demo

################################################
# Using docker to provide the dev container
# with lua runtime and luarocks package manager
################################################

DEV_IMAGE := openpolicyagent/kong-plugin-opa-devcontainer
DEV_TAG := latest

# Build Docker image with lua runtime and luarocks installed
# This image will be used to run luarocks and lucheck commands
docker-build-devcontainer:
	@echo "Building Docker image to compile and test the plugin..."
	docker build -t ${DEV_IMAGE}:${DEV_TAG} ./.devcontainer

# Delete Docker image built with lua runtime and luarocks
docker-rm-devcontainer:
	docker image rm ${DEV_IMAGE}:${DEV_TAG}

################################################
# Docker Integration Demo Image
################################################

DEMO_IMAGE := openpolicyagent/demo-kong-plugin-opa
DEMO_TAG = $(shell head -n 1 Dockerfile | cut -f 2 -d ':') # version of the Kong image this demo is built on

# Build Kong image with kong-plugin-opa installed
docker-build-demo:
	docker build -t ${DEMO_IMAGE}:${DEMO_TAG} .

# Publish the demo image to Docker Hub
docker-publish-demo: docker-build-demo
	docker push $(DEMO_IMAGE):${DEMO_TAG}

################################################
# Docker Command Helpers
################################################

# run command in a devcontainer
docker-run-devcontainer = docker run --rm -v "$(shell pwd)":/usr/kong -w /usr/kong ${DEV_IMAGE}:${DEV_TAG} $(1)

# helper command to run integration tests
docker-compose-test = docker-compose -f integration/docker-compose.yml -f integration/docker-compose.test.yml $(1)
