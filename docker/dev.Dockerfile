FROM golang:1.24.2-alpine3.21

ENV CGO_ENABLED=0 GOPATH= TZ=Europe/London

# curl is used for the e2e healthcheck.
RUN apk add --no-cache build-base curl git linux-headers tzdata \
 && GOBIN=/bin go install github.com/cespare/reflex@latest

COPY --from=codegolf/lang-uiua         ["/", "/langs/uiua/rootfs/"        ] # 19.8 MiB

COPY cmd/hash-langs ./cmd/hash-langs

RUN go run ./cmd/hash-langs/main.go

COPY run-lang.c ./

RUN gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

# reflex reruns a command when files change.
CMD reflex -sd none -r '\.(css|go|html|json|pem|svg|toml|txt)$' -R '_test\.go$' -- go run .
