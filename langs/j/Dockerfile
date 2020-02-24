FROM alpine:3.11 as builder

RUN apk add --no-cache build-base curl fts-dev

RUN curl -L https://github.com/jsoftware/jsource/archive/6a42dbc.tar.gz | tar xzf -

COPY patch /

RUN mv jsource-* jsource                                       \
 && patch -p0 <patch                                           \
 && rm -r /jsource/jlibrary/system/config/version.txt          \
          /jsource/jlibrary/system/defs                        \
          /jsource/jlibrary/system/main/socket.ijs             \
 && cd jsource/jsrc                                            \
 && echo '#define jbuilder  "code-golf.io"' >> jversion.h      \
 && echo '#define jlicense  "GPL3"'         >> jversion.h      \
 && echo '#define jplatform "linux"'        >> jversion.h      \
 && echo '#define jtype     "release"'      >> jversion.h      \
 && echo '#define jversion  "807"'          >> jversion.h      \
 && gcc -s -O2 -o jconsole jconsole.c jeload.c                 \
 && rm andjnative.c fnmatch.c jconsole.c jeload.c jep.c jfex.c \
 && gcc -s -O2 -fPIC -fwrapv -fno-strict-aliasing -shared -o libj.so *.c blis/*.c

COPY j.c /

RUN gcc -s -o j j.c

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1          /lib/
COPY --from=0 /jsource/jlibrary/bin/profile.ijs /usr/bin/profile.ijs
COPY --from=0 /jsource/jlibrary/system          /usr/system
COPY --from=0 /j                                \
              /jsource/jsrc/jconsole            \
              /jsource/jsrc/libj.so             /usr/bin/

ENTRYPOINT ["/usr/bin/j", "/tmp/code.ijs"]
