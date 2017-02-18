GO=$(shell which go)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o dist/ipservice_linux .
build-osx:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -o dist/ipservice_osx .
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -o dist/ipservice_win .

TAG ?= latest
docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o dist/ipservice
	docker build -t zensh/ipservice:$(TAG) .
run:
	docker run -d --name ipservice --rm -p 8080:8080 zensh/ipservice
	docker ps
