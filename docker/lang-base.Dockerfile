FROM alpine:3.19

WORKDIR /scratch

RUN mkdir dev etc tmp

RUN echo nobody:x:99:99::/: > etc/passwd

FROM scratch AS lang-base-no-proc

COPY --from=0 /scratch /

FROM scratch AS lang-base

COPY --from=0 /scratch     /
COPY --from=0 /scratch/dev /proc
