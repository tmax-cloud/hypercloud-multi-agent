FROM golang:1.16.3-buster as build

WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a ./cmd/agent-apiserver/agent-apiserver.go

FROM tmaxcloudck/hypercloud-multi-agent-base:b5.0.0.1
WORKDIR /go/src
COPY --from=build /go/src/agent-apiserver .
ENTRYPOINT ["/go/src/agent-apiserver"]
