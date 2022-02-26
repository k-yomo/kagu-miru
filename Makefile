.DEFAULT_GOAL := help

.PHONY: setup
setup: ## Setup local development environment
	command -v golangci-lint &> /dev/null || brew install golangci-lint
	command -v air &> /dev/null || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
	go mod download
	./scripts/install_go_tools.sh
	cd frontend && npm install

.PHONY: run
run: ## Run servers
	docker compose up pubsub elasticsearch & \
	(./scripts/create_local_pubsub_resources.sh && goreman -set-ports=false start) & \
	wait

.PHONY: run-item-index
run-item-index: ## Run item fetchers
	goreman -set-ports=false -f item_indexing.Procfile start

.PHONY: run-media-studio
run-media-studio:
	cd media_studio && npm run start

.PHONY: test ## Run tests
test:
	gotestsum -- -race -coverprofile=coverage.out $(TESTARGS) ./backend/...

.PHONY: test-cover
test-cover: test ## Run tests and show coverage
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: lint
lint: ## Run lint
	@golangci-lint run & \
	cd frontend && npm run lint & \
	wait

.PHONY: fmt
fmt: ## Format code
	goimports -w backend & \
	cd frontend && npm run fmt & \
	cd media_studio && npm run fmt & \
	wait

.PHONY: gen-graphql
gen-graphql: ## Generate graphql related code based on the schema
	go generate ./... & \
	cd frontend && npm run codegen & \
	cd media_studio && npm run codegen & \
	gqldoc -s defs/graphql/schema.graphql -o ./docs/graphql & \
	wait
	make fmt

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
