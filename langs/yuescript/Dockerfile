FROM debian:bookworm-slim AS builder

RUN apt-get update \
 && apt-get install -y build-essential curl

RUN curl -L https://github.com/pigpigyyy/Yuescript/archive/refs/tags/v0.24.1.tar.gz | tar xz --strip-components=1

ENV LDFLAGS=-static

RUN make \
 && make install

FROM codegolf/lang-base

COPY --from=0 /yue /usr/bin/

ENTRYPOINT ["yue"]

CMD ["-v"]
