

.PHONY: run
run:
	goreman -set-ports=false start

fmt:
	go fmt ./... & \
	goimports -w . & \
	cd frontend && npm run fmt & \
	cd terraform && terraform fmt -recursive & \
	wait

.PHONY: tf-symlink
tf-symlink:
	#cd ./terraform/dev && ln -sf ../shared/* .
	cd ./terraform/prod && ln -sf ../shared/* .

