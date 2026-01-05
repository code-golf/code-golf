FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes bison flex gcc make wget

ENV VER=10.0

ADD https://gitlab.winehq.org/wine/wine.git#wine-$VER /

RUN ./configure           \
    --disable-tests       \
    --disable-win16       \
    --enable-silent-rules \
    --enable-werror       \
    --enable-win64        \
    --prefix=/usr         \
    --without-alsa        \
    --without-coreaudio   \
    --without-freetype    \
    --without-gettext     \
    --without-gnutls      \
    --without-oss         \
    --without-pulse       \
    --without-x           \
 && make -j`nproc` install

RUN strip /usr/bin/wineserver \
          /usr/lib/wine/x86_64-unix/*.so

ADD --chmod=755 http://winetricks.org/winetricks /usr/bin

RUN ln -s /usr/bin/wine64 /usr/bin/wine

WORKDIR /Winedows

RUN WINEARCH='win64' WINEPREFIX='/Winedows/System64' winetricks

COPY vbscript.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/vbscript -s vbscript.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                  \
              /lib/x86_64-linux-gnu/libgcc_s.so.1              /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                      /lib64/
COPY --from=0 /usr/bin/vbscript                                \
              /usr/bin/wine                                    \
              /usr/bin/wineserver                              /usr/bin/
COPY --from=0 /usr/lib/locale                                  /usr/lib/locale
COPY --from=0 /usr/lib/wine/x86_64-unix/advapi32.dll.so        \
              /usr/lib/wine/x86_64-unix/combase.dll.so         \
              /usr/lib/wine/x86_64-unix/coml2.dll.so           \
              /usr/lib/wine/x86_64-unix/cscript.exe.so         \
              /usr/lib/wine/x86_64-unix/dnsapi.*so             \
              /usr/lib/wine/x86_64-unix/gdi32.dll.so           \
              /usr/lib/wine/x86_64-unix/hidclass.sys.so        \
              /usr/lib/wine/x86_64-unix/hidparse.sys.so        \
              /usr/lib/wine/x86_64-unix/iphlpapi.dll.so        \
              /usr/lib/wine/x86_64-unix/kernel32.dll.so        \
              /usr/lib/wine/x86_64-unix/kernelbase.dll.so      \
              /usr/lib/wine/x86_64-unix/mountmgr.*so           \
              /usr/lib/wine/x86_64-unix/msvcrt.dll.so          \
              /usr/lib/wine/x86_64-unix/ndis.sys.so            \
              /usr/lib/wine/x86_64-unix/nsi.dll.so             \
              /usr/lib/wine/x86_64-unix/nsiproxy.*so           \
              /usr/lib/wine/x86_64-unix/ntdll.*so              \
              /usr/lib/wine/x86_64-unix/ntoskrnl.exe.so        \
              /usr/lib/wine/x86_64-unix/ole32.dll.so           \
              /usr/lib/wine/x86_64-unix/oleaut32.dll.so        \
              /usr/lib/wine/x86_64-unix/plugplay.exe.so        \
              /usr/lib/wine/x86_64-unix/propsys.dll.so         \
              /usr/lib/wine/x86_64-unix/rpcrt4.dll.so          \
              /usr/lib/wine/x86_64-unix/sechost.dll.so         \
              /usr/lib/wine/x86_64-unix/services.exe.so        \
              /usr/lib/wine/x86_64-unix/setupapi.dll.so        \
              /usr/lib/wine/x86_64-unix/shcore.dll.so          \
              /usr/lib/wine/x86_64-unix/shell32.dll.so         \
              /usr/lib/wine/x86_64-unix/shlwapi.dll.so         \
              /usr/lib/wine/x86_64-unix/start.exe.so           \
              /usr/lib/wine/x86_64-unix/svchost.exe.so         \
              /usr/lib/wine/x86_64-unix/ucrtbase.dll.so        \
              /usr/lib/wine/x86_64-unix/user32.dll.so          \
              /usr/lib/wine/x86_64-unix/userenv.dll.so         \
              /usr/lib/wine/x86_64-unix/vbscript.dll.so        \
              /usr/lib/wine/x86_64-unix/wevtsvc.dll.so         \
              /usr/lib/wine/x86_64-unix/win32u.dll.so          \
              /usr/lib/wine/x86_64-unix/windowscodecs.dll.so   \
              /usr/lib/wine/x86_64-unix/wineboot.exe.so        \
              /usr/lib/wine/x86_64-unix/winebth.*so            \
              /usr/lib/wine/x86_64-unix/winebus.*so            \
              /usr/lib/wine/x86_64-unix/winedevice.exe.so      \
              /usr/lib/wine/x86_64-unix/winehid.sys.so         \
              /usr/lib/wine/x86_64-unix/winemenubuilder.exe.so \
              /usr/lib/wine/x86_64-unix/ws2_32.*so             /usr/lib/wine/x86_64-unix/
COPY --from=0 /usr/lib/wine/x86_64-windows/apisetschema.dll    /usr/lib/wine/x86_64-windows/
COPY --from=0 /usr/share/wine/wine.inf                         /usr/share/wine/
COPY --from=0 /usr/share/wine/nls/l_intl.nls                   \
              /usr/share/wine/nls/locale.nls                   \
              /usr/share/wine/nls/sortdefault.nls              /usr/share/wine/nls/
COPY --chown=nobody                                            \
     --from=0 /Winedows/System64                               /Winedows/System64

ENTRYPOINT ["vbscript"]

CMD ["--version"]
