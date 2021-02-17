FROM golang:1.16-buster

RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.34.1
