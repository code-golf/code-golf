FROM alpine:3.22 AS builder

ENV VER=1.2.1

RUN wget -O- https://github.com/dloscutoff/pip/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

FROM codegolf/lang-python

COPY --from=0 /errors.py    \
              /execution.py \
              /operators.py \
              /parsing.py   \
              /pip.py       \
              /ptypes.py    \
              /scanning.py  \
              /tokens.py    \
              /version.py   /usr/lib/

ENTRYPOINT ["python"]

CMD ["/usr/lib/pip.py", "--version"]
