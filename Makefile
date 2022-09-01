#!/bin/bash
# Makefile to wrap common docker and dev related tasks. Just type 'make' to get
# help.
#
# Requirements:
#  - Install Docker locally
#	 	Mac OS `brew cask install docker`
#    	Linux: `apt-get install docker`

REGISTRY = docker.io
REPOSITORY = registry
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
	@echo "Revision: $(REVISION)"
	@echo "Branch: $(BRANCH)"

build: ## Build the Registry container from current repo. Make sure to commit all changes beforehand
	docker build --build-arg REGISTRY_RELEASE=$(REVISION) -t registry:multi -t $(TAG) -t aptrust/$(TAG)-$(BRANCH) -t $(REGISTRY)/$(REPOSITORY)/$(TAG) -t $(REGISTRY)/$(REPOSITORY)/registry:$(REVISION)-$(BRANCH) -f Dockerfile.multi .

build-nc: ## Build the Registry container, no cached layers.
	docker build --no-cache --build-arg REGISTRY_RELEASE=$(REVISION) -t aptrust/$(TAG)-$(BRANCH) -t $(REGISTRY)/$(REPOSITORY)/$(TAG) .


up: ## Start containers for Registry service in docker-compose - currently not in service
	DOCKER_TAG_NAME=$(REVISION) docker-compose up

down: ## Stop containers for Registry for a docker compose deploymentn- currently not in service
	docker-compose down

run: ## Just run Registry in foreground
	docker run -p 8080:8080 $(TAG)

runshell: ## Run Registry container with interactive shell
	docker run -it --rm --env-file=.env $(REGISTRY)/$(REPOSITORY)/pharos:$(REVISION)-$(BRANCH) bash
	
registry_login: ## Log in to Docker Registry temporarily Docker Hub
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
	@echo "Pushing aptrust/registry:/$(REVISION)-$(BRANCH) for local deploy"
	@echo "Pushing aptrust/$(TAG)-$(PUSHBRANCH) info only for debugging on travis."

# Docker release - build, tag and push the container
release: build publish ## Make a release by building and publishing the `{version}` as `latest` tagged containers to Gitlab

push: ## Push the Docker image up to the registry
	docker push  $(REGISTRY)/$(REPOSITORY)/$(TAG)-$(BRANCH)

update-template: ## Update Cloudformation template with latest container version
	@echo "Overwriting container revision and branch from the CFN template to the CFN deployment YAML document."
	sed 's/registry:multi/registry:$(REVISION)-$(BRANCH)/g' cfn/cfn-registry-cluster.tmpl > cfn/cfn-registry-cluster.yml

clean: ## Clean the generated/compiles files


  
