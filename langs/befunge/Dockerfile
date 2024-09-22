FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl

RUN curl https://codeberg.org/catseye/Befunge-93/archive/rel_2_25.tar.gz \
  | tar xz --strip-components 1

RUN make

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /bin/bef                 /usr/bin/

COPY /befunge /usr/bin/

ENTRYPOINT ["befunge"]
