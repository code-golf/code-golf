FROM alpine:3.23 AS builder

ENV VER=0.3.0.1

RUN wget -qO /usr/lib/kamilalisp.jar https://github.com/kspalaiologos/kamilalisp/releases/download/v$VER/kamilalisp-$VER.jar

FROM codegolf/lang-java

COPY --from=0 /usr/lib/kamilalisp.jar /usr/lib/
