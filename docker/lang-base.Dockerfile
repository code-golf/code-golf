FROM alpine:3.15

WORKDIR /scratch

RUN mkdir dev etc proc tmp

RUN echo nobody:x:99:99::/: > etc/passwd

FROM scratch

COPY --from=0 /scratch /
