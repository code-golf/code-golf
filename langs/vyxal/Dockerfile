FROM alpine:3.23 AS builder

ENV VER=3.10.3

RUN wget -qO /usr/lib/vyxal.jar https://github.com/Vyxal/Vyxal/releases/download/v$VER/vyxal-$VER.jar

FROM codegolf/lang-java

COPY --from=0 /usr/lib/vyxal.jar /usr/lib/

# Delete default CDS archive to save space.
RUN rm -f /opt/java/lib/server/classes.jsa

FROM codegolf/lang-base

COPY --from=1 /lib /lib
COPY --from=1 /opt /opt
COPY --from=1 /usr /usr

CMD ["/opt/java/bin/java", "-jar", "/usr/lib/vyxal.jar", "--version"]
