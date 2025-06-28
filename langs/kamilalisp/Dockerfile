FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc musl-dev

ENV VER=0.3.0.1

RUN curl -#Lo kamilalisp.jar https://github.com/kspalaiologos/kamilalisp/releases/download/v$VER/kamilalisp-$VER.jar

COPY kamilalisp.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/kamilalisp -s -static kamilalisp.c

FROM codegolf/lang-java

COPY --from=0 /kamilalisp.jar     /
COPY --from=0 /usr/bin/kamilalisp /usr/bin/

ENTRYPOINT ["kamilalisp"]
