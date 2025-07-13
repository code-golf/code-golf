ARG VER=0.56

FROM alpine:3.22 AS builder

ARG VER

RUN apk add --no-cache bison curl flex g++ gobject-introspection-dev make

RUN curl -#L https://download.gnome.org/sources/vala/$VER/vala-$VER.18.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure                      \
    --disable-valadoc                \
    --prefix=/usr                    \
 && make -j`nproc` install           \
 && strip /usr/lib/libvala-$VER.so.0 \
          /usr/lib/vala-$VER/libvalaccodegen.so

COPY vala.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/vala -s vala.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/ld-musl-x86_64.so.1                   /lib/
COPY --from=0 /usr/bin/as                                \
              /usr/bin/cc                                \
              /usr/bin/gcc                               \
              /usr/bin/ld                                \
              /usr/bin/pkg-config                        \
              /usr/bin/vala                              \
              /usr/bin/valac                             /usr/bin/
COPY --from=0 /usr/include                               /usr/include
COPY --from=0 /usr/lib/Scrt1.o                           \
              /usr/lib/crti.o                            \
              /usr/lib/crtn.o                            \
              /usr/lib/libbfd-2.44.so                    \
              /usr/lib/libc.so                           \
              /usr/lib/libctf.so.0                       \
              /usr/lib/libffi.so.8                       \
              /usr/lib/libgcc_s.so                       \
              /usr/lib/libgcc_s.so.1                     \
              /usr/lib/libglib-2.0.so                    \
              /usr/lib/libglib-2.0.so.0                  \
              /usr/lib/libgmodule-2.0.so.0               \
              /usr/lib/libgmp.so.10                      \
              /usr/lib/libgobject-2.0.so                 \
              /usr/lib/libgobject-2.0.so.0               \
              /usr/lib/libintl.so                        \
              /usr/lib/libintl.so.8                      \
              /usr/lib/libisl.so.23                      \
              /usr/lib/libjansson.so.4                   \
              /usr/lib/libmpc.so.3                       \
              /usr/lib/libmpfr.so.6                      \
              /usr/lib/libpcre2-8.so.0                   \
              /usr/lib/libpkgconf.so.6                   \
              /usr/lib/libsframe.so.1                    \
              /usr/lib/libssp_nonshared.a                \
              /usr/lib/libvala-$VER.so.0                 \
              /usr/lib/libz.so.1                         \
              /usr/lib/libzstd.so.1                      /usr/lib/
COPY --from=0 /usr/lib/gcc/*/14.2.0/crtbeginS.o          \
              /usr/lib/gcc/*/14.2.0/crtendS.o            \
              /usr/lib/gcc/*/14.2.0/libgcc.a             /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/lib/glib-2.0/include                  /usr/lib/glib-2.0/include
COPY --from=0 /usr/lib/pkgconfig/glib-2.0.pc             \
              /usr/lib/pkgconfig/gobject-2.0.pc          \
              /usr/lib/pkgconfig/libffi.pc               \
              /usr/lib/pkgconfig/libpcre2-8.pc           /usr/lib/pkgconfig/
COPY --from=0 /usr/lib/vala-$VER/libvalaccodegen.so      /usr/lib/vala-$VER/
COPY --from=0 /usr/libexec/gcc/*/14.2.0/cc1              \
              /usr/libexec/gcc/*/14.2.0/liblto_plugin.so /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/share/gir-1.0/GLib-2.0.gir            \
              /usr/share/gir-1.0/GObject-2.0.gir         /usr/share/gir-1.0/
COPY --from=0 /usr/share/vala-$VER/vapi/glib-2.0.vapi    \
              /usr/share/vala-$VER/vapi/gobject-2.0.vapi /usr/share/vala-$VER/vapi/

ENTRYPOINT ["vala"]

CMD ["--version"]
