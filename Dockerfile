FROM golang:1.14.1-alpine

ENV CGO_ENABLED=0 GOPATH=

RUN apk add --no-cache build-base linux-headers

COPY go.* ./

RUN go mod download

COPY main.go run-lang.c ./
COPY cookie             ./cookie/
COPY github             ./github/
COPY hole               ./hole/
COPY pie                ./pie/
COPY pretty             ./pretty/
COPY routes             ./routes/

RUN go build -ldflags -s \
 && gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

RUN mkdir /empty

FROM scratch

COPY --from=codegolf/lang-bash       /       /langs/bash/rootfs/
COPY --from=codegolf/lang-brainfuck  /  /langs/brainfuck/rootfs/
COPY --from=codegolf/lang-c          /          /langs/c/rootfs/
COPY --from=codegolf/lang-haskell    /    /langs/haskell/rootfs/
COPY --from=codegolf/lang-j          /          /langs/j/rootfs/
COPY --from=codegolf/lang-javascript / /langs/javascript/rootfs/
COPY --from=codegolf/lang-julia      /      /langs/julia/rootfs/
COPY --from=codegolf/lang-lisp       /       /langs/lisp/rootfs/
COPY --from=codegolf/lang-lua        /        /langs/lua/rootfs/
COPY --from=codegolf/lang-nim        /        /langs/nim/rootfs/
COPY --from=codegolf/lang-perl       /       /langs/perl/rootfs/
COPY --from=codegolf/lang-php        /        /langs/php/rootfs/
COPY --from=codegolf/lang-python     /     /langs/python/rootfs/
COPY --from=codegolf/lang-raku       /       /langs/raku/rootfs/
COPY --from=codegolf/lang-ruby       /       /langs/ruby/rootfs/
COPY --from=codegolf/lang-rust       /       /langs/rust/rootfs/
COPY --from=codegolf/lang-swift      /      /langs/swift/rootfs/

COPY --from=0 /empty       /langs/bash/rootfs/proc/
COPY --from=0 /empty       /langs/bash/rootfs/tmp/
COPY --from=0 /empty  /langs/brainfuck/rootfs/proc/
COPY --from=0 /empty  /langs/brainfuck/rootfs/tmp/
COPY --from=0 /empty          /langs/c/rootfs/proc/
COPY --from=0 /empty          /langs/c/rootfs/tmp/
COPY --from=0 /empty    /langs/haskell/rootfs/proc/
COPY --from=0 /empty    /langs/haskell/rootfs/tmp/
COPY --from=0 /empty          /langs/j/rootfs/proc/
COPY --from=0 /empty          /langs/j/rootfs/tmp/
COPY --from=0 /empty /langs/javascript/rootfs/proc/
COPY --from=0 /empty /langs/javascript/rootfs/tmp/
COPY --from=0 /empty      /langs/julia/rootfs/proc/
COPY --from=0 /empty      /langs/julia/rootfs/tmp/
COPY --from=0 /empty       /langs/lisp/rootfs/proc/
COPY --from=0 /empty       /langs/lisp/rootfs/tmp/
COPY --from=0 /empty        /langs/lua/rootfs/proc/
COPY --from=0 /empty        /langs/lua/rootfs/tmp/
COPY --from=0 /empty        /langs/nim/rootfs/proc/
COPY --from=0 /empty        /langs/nim/rootfs/tmp/
COPY --from=0 /empty       /langs/perl/rootfs/proc/
COPY --from=0 /empty       /langs/perl/rootfs/tmp/
COPY --from=0 /empty        /langs/php/rootfs/proc/
COPY --from=0 /empty        /langs/php/rootfs/tmp/
COPY --from=0 /empty     /langs/python/rootfs/proc/
COPY --from=0 /empty     /langs/python/rootfs/tmp/
COPY --from=0 /empty       /langs/raku/rootfs/proc/
COPY --from=0 /empty       /langs/raku/rootfs/tmp/
COPY --from=0 /empty       /langs/ruby/rootfs/proc/
COPY --from=0 /empty       /langs/ruby/rootfs/tmp/
COPY --from=0 /empty       /langs/rust/rootfs/proc/
COPY --from=0 /empty       /langs/rust/rootfs/tmp/
COPY --from=0 /empty      /langs/swift/rootfs/proc/
COPY --from=0 /empty      /langs/swift/rootfs/tmp/

COPY --from=0 /go/code-golf                      /
COPY --from=0 /usr/bin/run-lang                  /usr/bin/
COPY          holes.toml                         /
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY          views                              /views/

CMD ["/code-golf"]
