FROM alpine:3.22 AS builder

ENV VER=1.2.0

RUN wget -O- https://github.com/dloscutoff/pip/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

FROM codegolf/lang-python

COPY --from=0 /*.py /

ENTRYPOINT ["python"]

CMD ["pip.py", "--version"]
