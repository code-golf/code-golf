FROM debian:bookworm-slim AS builder

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential curl

RUN curl -L https://github.com/pigpigyyy/Yuescript/archive/refs/tags/v0.24.1.tar.gz \
  | tar xz --strip-components 1

RUN LDFLAGS="-static" \
    make -j`nproc`    \
 && make install

FROM codegolf/lang-base

COPY --from=0 /bin     /bin
COPY --from=0 /lib     /lib
COPY --from=0 /lib64   /lib64
COPY --from=0 /usr/bin /usr/bin
COPY --from=0 /yue     /usr/bin/
COPY --from=0 /usr/lib /usr/lib

ENTRYPOINT ["yue"]

CMD ["-v"]
