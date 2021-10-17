
.PHONY: setup
setup:
	command -v golangci-lint &> /dev/null || brew install golangci-lint
	command -v air &> /dev/null || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
	go mod download
	./scripts/install_go_tools.sh
	cd frontend && npm install

.PHONY: run
run:
	docker compose up elasticsearch & \
	docker compose up pubsub & \
	./scripts/create_local_pubsub_resources.sh & \
	goreman -set-ports=false start & \
	wait

.PHONY: test
test:
	gotestsum -- -race -coverprofile=coverage.out $(TESTARGS) ./backend/...

.PHONY: test-cover
test-cover: testacc
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: lint
lint:
	@golangci-lint run & \
	cd frontend && npm run lint & \
	wait

.PHONY: fmt
fmt:
	go fmt ./... & \
	goimports -w . & \
	cd frontend && npm run fmt & \
	wait

.PHONY: gen-graphql
gen-graphql:
	go generate ./... & \
	cd frontend && npm run codegen & \
	gqldoc -s defs/graphql/schema.graphql -o ./docs/graphql & \
	wait
	make fmt
