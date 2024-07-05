SHELL:=/bin/bash

.PHONY: docker-network
docker-network:
	docker network create dating-app-network


.PHONY: up
up:
	docker-compose up -d

.PHONY: build
build:
	docker-compose build

.PHONY: build-and-up
build-and-up: build up
