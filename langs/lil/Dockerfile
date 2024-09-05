FROM alpine:3.20 AS builder

RUN apk add --no-cache bash build-base curl sdl2-dev

RUN curl -L https://github.com/JohnEarnest/Decker/archive/refs/tags/v1.48.tar.gz \
  | tar xz --strip-components 1

RUN make lilt \
 && make install

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /c/build/lilt            /usr/bin/

COPY /lilwrapper /usr/bin/

ENTRYPOINT ["lilwrapper"]

CMD ["--version"]
