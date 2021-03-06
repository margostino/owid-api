#!make
SHELL = /bin/sh
.DEFAULT: help

-include .env .env.local .env.*.local

# Defaults
DOCKER_COMPOSE = USERID=$(shell id -u):$(shell id -g) docker-compose ${compose-files}
ALL_ENVS := local ci
env ?= local

.PHONY: help
help:
	@echo "Lumos pipeline"
	@echo ""
	@echo "Usage:"
	@echo "  vercel.dev                     - Replicate the Vercel deployment environment locally, allowing you to test your Serverless Functions, without requiring you to deploy each time a change is made"
	@echo "  run                            - Start main runner for local testing"
	@echo ""
	@echo "Project-level environment variables are set in .env file:"
	@echo "  VERCEL=1"
	@echo "  VERCEL_ENV=development"
	@echo ""
	@echo "Note: Store protected environment variables in .env.local or .env.*.local"
	@echo ""

.PHONY: vercel.pull
vercel.pull:
	vercel pull

.PHONY: vercel.dev
vercel.dev:
	make v.stop
	vercel dev

.PHONY: run
run:
	go run runner/main.go

.PHONY: v.stop
v.stop:
	./bin/vercel-stop

.PHONY: schema
schema:
	go run ./runner schema-gen

.PHONY: server
server:
	go get github.com/99designs/gqlgen@v0.17.9
	go run github.com/99designs/gqlgen generate

.PHONY: clean.metadata
clean.metadata:
	rm ./metadata/*

.PHONY: recovery
recovery:
	git checkout -m graph/generated/generated.go
	git checkout -m graph/model/models_gen.go

.PHONY: template.copy
template.copy:
	sudo cp ./templates/custom_resolver.gotpl /Users/martin.dagostino/go/pkg/mod/github.com/99designs/gqlgen@v0.17.9/plugin/resolvergen/
