FROM codegolf/lang-base AS builder

ENV VER=1.13.220

ADD --unpack https://github.com/babashka/babashka/releases/download/v$VER/babashka-$VER-linux-amd64-static.tar.gz /usr/bin

ENTRYPOINT ["bb"]

CMD ["--version"]
