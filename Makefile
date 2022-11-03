PATH := $(CURDIR)/bin:$(PATH)

DOCKER_COMPOSE := $(or $(DOCKER_COMPOSE),$(DOCKER_COMPOSE),docker-compose)

run:
	$(DOCKER_COMPOSE) up
