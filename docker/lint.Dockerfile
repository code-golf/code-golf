FROM golang:1.18.0-bullseye

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.0
