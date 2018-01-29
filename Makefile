BEAT_NAME=openconfigbeat
BEAT_DIR=github.com/aristanetworks
BEAT_PATH=$(BEAT_DIR)/$(BEAT_NAME)
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
GOIMPORTS_LOCAL_PREFIX=github.com/aristanetworks
PREFIX ?= .
DOCKER = docker
ELASTICSEARCH_VERSION = 6.1.2
ELASTICSEARCH_HOST ?= 127.0.0.1
DOCKER_IMAGE = docker.elastic.co/elasticsearch/elasticsearch:$(ELASTICSEARCH_VERSION)
DOCKER_CONTAINER = openconfigbeat-elasticsearch
GO = go
GOPKGVERSION := $(shell git describe --tags --match "[0-9]*" --abbrev=7 HEAD)
ifndef GOPKGVERSION
   $(error unable to determine git version)
endif
GOBUILD_FLAGS ?= -ldflags "-s -w -X github.com/aristanetworks/openconfigbeat/cmd.Version=$(GOPKGVERSION)"

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile

.PHONY: collect

.PHONY: update-deps
update-deps:
	glide update && ./clean_vendor.sh

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:

.PHONY: docker-start
docker-start:
	@$(MAKE) docker-stop || true
	$(DOCKER) run --name $(DOCKER_CONTAINER) -d -p $(ELASTICSEARCH_HOST):9200:9200 $(DOCKER_IMAGE)
	echo "Waiting for elasticsearch to be reachable..." && time sh -c "until curl -sf http://$(ELASTICSEARCH_HOST):9200; do sleep 1; done"

.PHONY: docker-stop
docker-stop:
	$(DOCKER) stop $(DOCKER_CONTAINER) && $(DOCKER) rm $(DOCKER_CONTAINER)

.PHONY: beater-test
beater-test: $(BEAT_NAME)
	$(MAKE) docker-start
	$(MAKE) testsuite
	$(MAKE) docker-stop
