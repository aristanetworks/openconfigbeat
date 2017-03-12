BEAT_NAME=openconfigbeat
BEAT_DIR=github.com/aristanetworks
BEAT_PATH=$(BEAT_DIR)/$(BEAT_NAME)
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile

.PHONY: collect

.PHONY: update-deps
update-deps:
	glide update --no-recursive

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:
