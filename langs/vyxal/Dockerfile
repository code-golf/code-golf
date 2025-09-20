FROM alpine:3.22 AS builder

ENV VER=3.10.1

RUN wget -qO vyxal.jar https://github.com/Vyxal/Vyxal/releases/download/v$VER/vyxal-$VER.jar

FROM codegolf/lang-java

COPY --from=0 /vyxal.jar /

CMD ["/opt/java/bin/java", "-jar", "vyxal.jar", "-h"]
