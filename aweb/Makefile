.PHONY: docker
docker:
	@rm aweb || true
	@go mod tidy
	@GOOS=linux GOARCH=amd64 go build -tags=k8s -o aweb .
	@docker rmi -f pluckhuang/aweb:v0.0.1
	@docker build -t pluckhuang/aweb:v0.0.1 .


.PHONY: build
build:
	@rm aweb || true
	@go mod tidy
	@go build -tags=dev -o aweb .



clean:
	@rm aweb || true