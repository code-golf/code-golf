FROM debian:bookworm-slim AS builder

ENV DEBIAN_FRONTEND=noninteractive LDFLAGS=-static

RUN apt-get update \
 && apt-get install -y build-essential curl

RUN curl -L https://github.com/pigpigyyy/Yuescript/archive/refs/tags/v0.25.5.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` install

FROM codegolf/lang-base

COPY --from=0 /yue /usr/bin/

ENTRYPOINT ["yue"]

CMD ["-v"]
