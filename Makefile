# note: call scripts from /scripts

IMG ?= tmaxcloudck/hypercloud-multi-agent:b5.0.0.1

.PHONY: build
build: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' ./cmd/agent-apiserver/agent-apiserver.go

.PHONY: docker-build
docker-build: 
	docker build -f ./build/Dockerfile -t $(IMG) .

.PHONY: docker-build
docker-push: 
	docker push ${IMG}
