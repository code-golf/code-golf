FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc linux-headers make musl-dev

ENV HB_BUILD_CONTRIBS='no' HB_BUILD_STRIP='all' HB_USER_LDFLAGS='-static' VER=3.0.0

RUN curl -#L https://downloads.sourceforge.net/harbour-project/harbour-$VER.tar.gz \
  | tar xz --strip-components 1

RUN curl -#L https://github.com/harbour/core/archive/master.tar.gz \
  | tar xz --directory /tmp --strip-components 1
# FIXME File "arc4.c" still uses deprecated system call `sysctl`. The error
# during compilation can be suppressed, but not the UR during linking.
# See: https://www.man7.org/linux/man-pages/man2/sysctl.2.html and
#      https://github.com/harbour/core/blob/cfb7bdc/src/rtl/arc4.c#L57-L58
RUN cp /tmp/src/rtl/arc4.c src/rtl/arc4.c \
 && cp /tmp/src/rtl/arc4.h include/hbarc4.h

RUN make -j`nproc` install

WORKDIR /tmp

ENV HB_INSTALL_PREFIX='/usr'
# FIXME This is truly awful. The binary "hbmk2" dates back over ten years and can
# execute scripts in Harbour, but it was still in an experimental phase. To put it
# bluntly, it's completely useless, as it didn't work then and doesn't work now.
# Therefore, a second, more recent version of "hbmk2" is being installed here from
# the master branch of the GitHub repository, which thankfully does work and will
# serve as a replacement for the first installed version.
RUN make -j`nproc` install

COPY harbour.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/harbour -s -static /harbour.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/harbour           /usr/bin/
COPY --from=0 /usr/bin/hbmk2             /usr/local/bin/
COPY --from=0 /usr/local/include/harbour /usr/local/include/harbour
# FIXME Here, it was possible to omit a file using the version control hack in the
# "build-langs" script. That's what it's for, after all, especially after all that
# fiddling with overwriting binary files. Unfortunately, the new file "hbmk2" has
# a different version number.
ENTRYPOINT ["harbour"]

# CMD ["--version"]
