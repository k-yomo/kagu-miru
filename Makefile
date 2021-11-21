
.PHONY: setup
setup:
	command -v golangci-lint &> /dev/null || brew install golangci-lint
	command -v air &> /dev/null || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
	go mod download
	./scripts/install_go_tools.sh
	cd frontend && npm install

.PHONY: run
run:
	docker compose up pubsub elasticsearch & \
	(./scripts/create_local_pubsub_resources.sh && goreman -set-ports=false start) & \
	wait

.PHONY: run-item-fetcher
run-item-fetcher:
	goreman -set-ports=false -f item_fetcher.Procfile start

.PHONY: run-media-studio
run-media-studio:
	cd media_studio && npm run start

.PHONY: test
test:
	gotestsum -- -race -coverprofile=coverage.out $(TESTARGS) ./backend/...

.PHONY: test-cover
test-cover: test
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: lint
lint:
	@golangci-lint run & \
	cd frontend && npm run lint & \
	wait

.PHONY: fmt
fmt:
	goimports -w backend & \
	cd frontend && npm run fmt & \
	wait

.PHONY: gen-graphql
gen-graphql:
	go generate ./... & \
	cd frontend && npm run codegen & \
	cd media_studio && npm run codegen & \
	gqldoc -s defs/graphql/schema.graphql -o ./docs/graphql & \
	wait
	make fmt
