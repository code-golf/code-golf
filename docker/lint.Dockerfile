FROM golang:1.16beta1-buster

RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.34.1
