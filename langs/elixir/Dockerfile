# FIXME Using alpine causes a 5s timeout with no net, musl DNS???
# <1s docker run --rm --net none elixir:slim
# >5s docker run --rm --net none elixir:alpine
FROM elixir:1.13.4-slim as builder

FROM codegolf/lang-base

COPY --from=0 /bin                  /bin
COPY --from=0 /lib/x86_64-linux-gnu /lib/x86_64-linux-gnu
COPY --from=0 /lib64                /lib64
COPY --from=0 /usr                  /usr

ENTRYPOINT ["elixir"]

CMD ["--short-version"]
