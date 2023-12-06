# This file is just a porcelaine for devbox commands.
# Devbox is a tool to manage development environment dependencies.
# See https://www.jetpack.io/devbox for more information.

# You can install devbox with the following command:
# make bootstrap
# It also installs direnv to automatically load devbox environment.

# If you don't want to use devbox, you can use a containerized version of devbox
# with all the necessary dependencies. To do so, you need to install docker (or podman, etc)
# and run the targets with DOCKER=1 prefix. For example:
# make DOCKER=1 build

bootstrap:
	curl -fsSL https://get.jetpack.io/devbox | bash
	curl -sfL https://direnv.net/install.sh | bash

.PHONY: format
ifndef DOCKER
format: check_devbox
	devbox run format
else
format: $(DOCKER_BUILDER)
	$(call docker_builder_devbox,run format)
endif

.PHONY: build
ifndef DOCKER
build: check_devbox
	devbox run build
else
build: $(DOCKER_BUILDER)
	$(call docker_builder_devbox,run build)
endif

.PHONY: generate
ifndef DOCKER
generate: check_devbox
	devbox run generate
else
generate: $(DOCKER_BUILDER)
	$(call docker_builder_devbox,run generate)
endif

check_%:
	@command -v $* >/dev/null || (echo "missing required tool $*" ; false)

CMD_DOCKER ?= docker
DOCKER_BUILDER ?= parca-dev/testdata-builder

.PHONY: $(DOCKER_BUILDER)
$(DOCKER_BUILDER): Dockerfile | check_$(CMD_DOCKER)
	$(CMD_DOCKER) build -t $(DOCKER_BUILDER):latest .

define docker_builder_devbox
	$(CMD_DOCKER) run --rm \
	-w /code \
	-v $(PWD):/code \
	--entrypoint devbox \
	$(DOCKER_BUILDER) $(1)
endef
