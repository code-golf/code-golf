FROM golang:1.20

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
