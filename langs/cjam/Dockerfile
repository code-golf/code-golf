FROM alpine:3.23 AS builder

ENV VER=0.6.5

RUN wget -qO /usr/lib/cjam.jar https://downloads.sourceforge.net/cjam/cjam-$VER.jar

FROM codegolf/lang-java

COPY --from=0 /usr/lib/cjam.jar /usr/lib/

# Delete default CDS archive to save space.
RUN rm -f /opt/java/lib/server/classes.jsa

FROM codegolf/lang-base

COPY --from=1 /lib /lib
COPY --from=1 /opt /opt
COPY --from=1 /usr /usr
