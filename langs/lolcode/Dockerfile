FROM alpine:3.13 as builder

RUN mkdir /empty

RUN apk add --no-cache curl build-base unzip python2 cmake

RUN curl -L https://github.com/justinmeza/lci/archive/master.zip --output lci.zip

RUN unzip lci.zip

RUN cd lci-master \
  && python2.7 install.py

# FROM scratch

# COPY --from=0 /empty             /proc
# COPY --from=0 /empty             /tmp
# COPY --from=0 /bin/sh            /bin/
# COPY --from=0 /usr/local/bin/lci /usr/bin/

ENTRYPOINT [ "/bin/sh", "-c", "/usr/local/bin/lci -v 2>&1" ]