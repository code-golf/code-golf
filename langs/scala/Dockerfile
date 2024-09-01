FROM alpine:3.20 AS builder

RUN apk add --no-cache bash curl git openjdk21

RUN wget https://github.com/lampepfl/dotty/releases/download/3.5.0/scala3-3.5.0.zip && \
    unzip scala3-3.5.0.zip -d /usr/local/scala && \
    rm scala3-3.5.0.zip

RUN curl -L https://github.com/VirtusLab/scala-cli/archive/refs/tags/v1.5.0.tar.gz | tar xz

RUN mkdir -p /root/.cache/coursier/v1/https/github.com/coursier/coursier/releases/download/v2.1.10

RUN cd /root/.cache/coursier/v1/https/github.com/coursier/coursier/releases/download/v2.1.10 \
 && wget https://github.com/coursier/coursier/releases/download/v2.1.10/cs-x86_64-pc-linux.gz \
 && mkdir -p cs-x86_64-pc-linux.gz \
 && gunzip -c cs-x86_64-pc-linux.gz > cs-x86_64-pc-linux.gz/cs-x86_64-pc-linux \
 && chmod +x cs-x86_64-pc-linux.gz/cs-x86_64-pc-linux

ENV SCALA_HOME=/usr/local/scala
ENV PATH=/usr/local/scala/bin:/usr/local/bin:$PATH

RUN cd scala-cli-1.5.0 \
 && ./mill 'cli[]'.nativeImage

FROM codegolf/lang-base

COPY --from=0 / /

ENTRYPOINT ["scala"]

CMD ["-version"]
