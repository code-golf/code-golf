FROM alpine:3.22 AS builder

ENV VER=3.10.3

RUN wget -qO vyxal.jar https://github.com/Vyxal/Vyxal/releases/download/v$VER/vyxal-$VER.jar

FROM codegolf/lang-java

COPY --from=0 /vyxal.jar /

ENTRYPOINT ["/opt/java/bin/java"]

CMD ["-jar", "vyxal.jar", "-h"]
