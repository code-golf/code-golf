FROM alpine:3.24 AS builder

RUN apk add --no-cache unzip gcc musl-dev

ENV VER=4.33.0

WORKDIR /

RUN wget -cO /lean.zip https://github.com/leanprover/lean4/releases/download/v$VER/lean-$VER-linux.zip \
 && unzip lean.zip \
 && mv lean-$VER-linux lean

COPY leanwrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/lean -s -static leanwrapper.c

FROM codegolf/lang-base

# copy everything verbatim for now
COPY --from=0 /lean         /usr/local/lean/
COPY --from=0 /usr/bin/lean /usr/bin/

ENTRYPOINT ["/usr/bin/lean"]

CMD ["--version"]
