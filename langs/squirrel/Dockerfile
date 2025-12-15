FROM alpine:3.23 AS builder

RUN apk add --no-cache clang cmake curl make patch

ENV VER=3.2

RUN curl -#L https://github.com/albertodemichelis/squirrel/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

COPY squirrel.patch /

RUN patch -p0 < squirrel.patch

WORKDIR /squirrel

RUN cmake ..                             \
    -DCMAKE_BUILD_TYPE='MinSizeRel'      \
    -DCMAKE_POLICY_VERSION_MINIMUM='3.5' \
 && make install

COPY squirrel.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/squirrel -s /squirrel.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1      /lib/
COPY --from=0 /usr/bin/squirrel        /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/
COPY --from=0 /usr/local/bin/sq_static /usr/local/bin/

ENTRYPOINT ["squirrel"]

CMD ["--version"]
