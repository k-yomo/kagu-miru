
.PHONY: setup
setup:
	brew list golangci-lint &>/dev/null || brew install golangci-lint
	command -v air &> /dev/null || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
	go mod download
	./scripts/install_go_tools.sh
	cd frontend && npm install

.PHONY: run
run:
	docker compose up elasticsearch & \
	goreman -set-ports=false start & \
	wait

lint:
	@golangci-lint run & \
	cd frontend && npm run lint & \
	wait

fmt:
	go fmt ./... & \
	goimports -w . & \
	cd frontend && npm run fmt & \
	wait

gen-graphql:
	go generate ./... & \
	cd frontend && npm run codegen & \
	gqldoc -s defs/graphql/schema.graphql -o ./docs/graphql & \
	wait
