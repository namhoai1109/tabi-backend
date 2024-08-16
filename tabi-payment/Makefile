.PHONY: build clean deploy

HOST ?= localhost:3003

## ===== Local =====

depends: ## Install & build dependencies
	go get ./...
	go build ./...
	go mod tidy

mod.clean:
	go clean -cache
	go clean -modcache

mod: ## Update dependencies
	go mod tidy && go mod vendor

provision: depends ## Provision dev environment
	sudo docker compose up -d
	scripts/waitdb.sh
	@$(MAKE) migrate

docker.run:
	sudo docker compose up -d

docker.run.hieu-macos:
	docker compose up -d

start: ## Bring up the server on dev environment
	go run cmd/api/main.go

migrate: ## Run database migrations
	go run cmd/migration/main.go

migrate.undo: ## Undo the last database migration
	go run cmd/migration/main.go --down

build.linux: clean ## Build the server binary file for Linux host
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 sh scripts/build.sh

clean: 
	rm -rf ./server ./*.out
	rm -rf .serverless 0

specs: ## Generate swagger specs
	HOST=$(HOST) sh scripts/specs-gen.sh

build.diagram: ## Build database diagram
	dbdocs build diagram/db.dbml

## ===== AWS Lambda =====

deploy: build ## Deploy the server to AWS Lambda
	sls deploy

deploy.func: ## Deploy the functions to AWS Lambda
	sh scripts/sls-func.sh deploy

invoke.migration: ## Invoke AWS Lambda Migration
	sh scripts/sls-func.sh invoke --function Migration