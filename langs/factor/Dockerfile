FROM debian:bookworm-slim AS builder

RUN apt-get update && apt-get install -y curl

RUN curl https://downloads.factorcode.org/releases/0.99/factor-linux-x86-64-0.99.tar.gz \
  | tar xz

FROM codegolf/lang-base

COPY --from=0 / /

ENTRYPOINT ["/factor/factor"]

CMD ["--version"]
