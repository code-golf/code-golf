FROM debian:bookworm-slim AS builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
 && apt-get install -y build-essential curl

RUN curl https://www.lua.org/ftp/lua-5.4.7.tar.gz | tar xz --strip-components=1

RUN make clean \
 && CFLAGS='-fPIC' make linux \
 && make install

RUN curl -L https://github.com/pigpigyyy/Yuescript/archive/refs/tags/v0.24.1.tar.gz | tar xz --strip-components=1

RUN make shared \
 && make install

FROM codegolf/lang-base

COPY --from=0 / /

ENTRYPOINT ["yue"]

CMD ["-v"]
