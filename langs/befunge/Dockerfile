FROM debian:bullseye-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes cmake curl

ENV LDFLAGS='-static' VER=b7bf640

RUN curl -#L https://github.com/VorpalBlade/cfunge/tarball/$VER \
  | tar xz --strip-components 1

WORKDIR /build

RUN cmake .. && make install \
 && strip /usr/local/bin/cfunge

COPY befunge.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/befunge -s -static /befunge.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/befunge /usr/local/bin/cfunge /usr/bin/

ENTRYPOINT ["befunge"]

CMD ["--version"]
