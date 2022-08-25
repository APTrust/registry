#!/bin/bash
# Makefile to wrap common docker and dev related tasks. Just type 'make' to get
# help.
#
# Requirements:
#  - Install Docker locally
#	 	Mac OS `brew cask install docker`
#    	Linux: `apt-get install docker`

REGISTRY = hub.docker.com
REPOSITORY = container-registry
NAME=$(shell basename $(CURDIR))
REVISION=$(shell git log -1 --pretty=%h)
REVISION=$(shell git rev-parse --short=7 HEAD)
BRANCH = $(subst /,_,$(shell git rev-parse --abbrev-ref HEAD))
PUSHBRANCH = $(subst /,_,$(TRAVIS_BRANCH))
TAG = $(NAME):$(REVISION)

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help build publish clean release test

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

revision: ## Show me the git hash
	@echo $(REVISION)
	@echo $(BRANCH)

build: ## Build the Registry container from current repo. Make sure to commit all changes beforehand
	docker build --build-arg REGISTRY_RELEASE=$(REVISION) -t registry:multi -t $(TAG) -t aptrust/$(TAG)-$(BRANCH) -t $(REGISTRY)/$(REPOSITORY)/$(TAG) -t $(REGISTRY)/$(REPOSITORY)/registry:$(REVISION)-$(BRANCH) -f Dockerfile.multi .
#	docker build --build-arg PHAROS_RELEASE=${REVISION} -t nginx-proxy-pharos:latest -t $(REGISTRY)/$(REPOSITORY)/nginx-proxy-pharos:$(REVISION) -t $(REGISTRY)/$(REPOSITORY)/nginx-proxy-pharos:$(REVISION)-$(BRANCH) -t aptrust/nginx-proxy-pharos -f Dockerfile.nginx .

build-nc: ## Build the Registry container, no cached layers.
	docker build --no-cache --build-arg REGISTRY_RELEASE=$(REVISION) -t aptrust/$(TAG) -t $(REGISTRY)/$(REPOSITORY)/$(TAG) .
#	docker build --build-arg PHAROS_RELEASE=${REVISION} -t $(REGISTRY)/$(REPOSITORY)/nginx-proxy-pharos:$(REVISION) -t $(REGISTRY)/$(REPOSITORY)/nginx-proxy-pharos -t aptrust/nginx-proxy-pharos -f Dockerfile.nginx .

up: ## Start containers for Pharos, Postgresql, Nginx
	DOCKER_TAG_NAME=$(REVISION) docker-compose up

down: ## Stop containers for Pharos, Postgresql, Nginx
	docker-compose down

run: ## Just run Registry in foreground
	docker run -p 8080:8080 $(TAG)

runshell: ## Run Pharos container with interactive shell
	docker run -it --rm --env-file=.env $(REGISTRY)/$(REPOSITORY)/pharos:$(REVISION)-$(BRANCH) bash

#runconsole: ## Run Rails Console
#	docker run -it --rm --env-file=.env $(REGISTRY)/$(REPOSITORY)/pharos:$(REVISION)-$(BRANCH) /bin/bash -c "export TERM=dumb && rails c"

runcmd: ## Start Pharos container, run command and exit.
	docker run -it --rm --env-file=.env $(REGISTRY)/$(REPOSITORY)/registry:$(REVISION)-$(BRANCH) $(filter-out $@, $(MAKECMDGOALS))

%:
	@true



registry_login: ## Log in to Docker Registry
	# GITLAB
	#docker login $(REGISTRY)
	# Docker Hub
	docker login docker.io
#	docker push aptrust/registry:$(REVISION)-$(BRANCH)

publish: registry_login
	# GITLAB
#	docker login $(REGISTRY)
#	docker push $(REGISTRY)/$(REPOSITORY)/pharos
#	docker push $(REGISTRY)/$(REPOSITORY)/registry:$(REVISION)-$(BRANCH)
#	docker push $(REGISTRY)/$(REPOSITORY)/nginx-proxy-pharos:$(REVISION)-$(BRANCH)
#	#docker build --build-arg PHAROS_RELEASE=${REVISION} -t $(REGISTRY)/$(REPOSITORY)/nginx-proxy-pharos -t aptrust/nginx-proxy-pharos -f Dockerfile.nginx .
	# Docker Hub
	docker login docker.io
	docker push aptrust/registry:$(REVISION)-$(BRANCH)

publish-ci:
	@echo $(DOCKER_PWD) | docker login -u $(DOCKER_USER) --password-stdin $(REGISTRY)
	docker tag  $(REGISTRY)/$(REPOSITORY)/pharos:$(REVISION) $(REGISTRY)/$(REPOSITORY)/pharos:$(REVISION)-$(PUSHBRANCH)
	docker tag  $(REGISTRY)/$(REPOSITORY)/nginx-proxy-pharos:$(REVISION) $(REGISTRY)/$(REPOSITORY)/nginx-proxy-pharos:$(REVISION)-$(PUSHBRANCH)
	#docker push $(REGISTRY)/$(REPOSITORY)/pharos
	docker push $(REGISTRY)/$(REPOSITORY)/pharos:$(REVISION)-$(PUSHBRANCH)
	docker push $(REGISTRY)/$(REPOSITORY)/nginx-proxy-pharos:$(REVISION)-$(PUSHBRANCH)
	# Docker Hub
	#docker login docker.io
	#docker push aptrust/pharos

# Docker release - build, tag and push the container
release: build publish ## Make a release by building and publishing the `{version}` as `latest` tagged containers to Gitlab

push: ## Push the Docker image up to the registry
	docker push  $(REGISTRY)/$(REPOSITORY)/$(TAG)

clean: ## Clean the generated/compiles files
