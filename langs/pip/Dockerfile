FROM alpine:3.20 AS builder

RUN apk add --no-cache curl

RUN curl -L https://github.com/dloscutoff/pip/archive/refs/tags/v1.2.0.tar.gz \
  | tar -xz --strip-components 1

FROM codegolf/lang-python

COPY --from=0 *.py /

FROM codegolf/lang-base

COPY --from=1 . /

CMD ["python", "pip.py", "--version"]
