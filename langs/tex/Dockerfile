FROM reitzig/texlive-minimal:2025.5 AS builder

RUN apk add --no-cache dash gcc musl-dev

ENV VER=0.2.1

# Download dvi-to-text.
RUN wget -O tex.tar.xz https://github.com/jared-hughes/dvi-to-text/releases/download/v$VER/dvi-to-text_v${VER}_x86_64-unknown-linux-musl.tar.xz \
 && tar xJf tex.tar.xz --directory /

WORKDIR /usr/local/texlive/2025/texmf-dist

# Add a full 8-bit font to allow emitting non-ascii bytes.
COPY octet.mf fonts/source/public/cm

# Update kpathsea file database.
RUN ls -R . > ls-R                                           \
  # Remove all ligatures.
 && rm -r fonts/tfm/public/cm/*.tfm                          \
 && sed -i 's/ligs:=./ligs:=0/g' fonts/source/public/cm/*.mf \
  # Load octet font.
 && echo "\font\octet=octet" >> tex/plain/base/plain.tex     \
  # Apply changes.
 && mktexfmt tex.fmt

# Prune unneeded stuff.
RUN bash -c "rm -r ../texmf-var/web2c/{lua*,pdftex,metafont}                              \
                   ../bin/x86_64-linuxmusl/{lua*,pdf*,*dvi*,mf-nowin,bibtex}              \
                   {scripts,bibtex,dvips,tex,fonts,xdvi,metafont,makeindex,mft,dvipdfmx}" \
 && rm -r ../texmf-var/fonts ../tlpkg

COPY tex.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/tex -s /tex.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash             /bin/sh
COPY --from=0 /lib/ld-musl-*.so.1       /lib/
COPY --from=0 /dvi-to-text /usr/bin/tex /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6 \
              /usr/lib/libreadline.so.8 /usr/lib/
COPY --from=0 /usr/local/bin            /usr/local/bin
COPY --from=0 /usr/local/texlive/2025   /usr/local/texlive/2025

ENTRYPOINT ["tex"]

CMD ["--version"]
