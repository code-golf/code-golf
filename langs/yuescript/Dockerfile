FROM debian:bookworm-slim AS builder

RUN apt-get update \
 && apt-get install -y build-essential curl

ENV LDFLAGS='-static'

RUN curl https://lua.org/ftp/lua-5.4.7.tar.gz | tar xz --strip-components=1

RUN make

RUN curl -L https://github.com/pigpigyyy/Yuescript/archive/refs/tags/v0.24.1.tar.gz | tar xz --strip-components=1

RUN make \
 && make install

FROM codegolf/lang-base

COPY --from=0 /usr/local/bin/yue /usr/bin/

ENTRYPOINT ["yue"]

CMD ["-v"]
