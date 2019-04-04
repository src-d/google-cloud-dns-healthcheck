# Package configuration
PROJECT = google-cloud-dns-healthcheck
COMMANDS = cmd/google-cloud-dns-healthcheck
DOCKERFILES = Dockerfile:$(PROJECT)

# Including ci Makefile
CI_REPOSITORY ?= https://github.com/src-d/ci.git
CI_PATH ?= $(shell pwd)/.ci
CI_VERSION ?= v1

MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	git clone --quiet --branch $(CI_VERSION) --depth 1 $(CI_REPOSITORY) $(CI_PATH);

-include $(MAKEFILE)

GO_BUILD_ENV = CGO_ENABLED=0
