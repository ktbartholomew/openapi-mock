docker-build:
	docker build -t quay.io/ktbartholomew/openapi-mock .
build:
	go build -a -o openapi-mock .
test:
	go test -v ./...