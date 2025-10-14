FROM alpine:3.22 AS builder

RUN apk add --no-cache gcc musl-dev

ENV VER=0.3.0.1

RUN wget -qO /usr/lib/kamilalisp.jar https://github.com/kspalaiologos/kamilalisp/releases/download/v$VER/kamilalisp-$VER.jar

COPY kamilalisp.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/kamilalisp -s -static kamilalisp.c

FROM codegolf/lang-java

COPY --from=0 /usr/bin/kamilalisp     /usr/bin/
COPY --from=0 /usr/lib/kamilalisp.jar /usr/lib/

ENTRYPOINT ["kamilalisp"]
