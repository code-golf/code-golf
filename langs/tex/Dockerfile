FROM reitzig/texlive-base:2022.8 AS builder

# Download dvi-to-text
ENV VERSION=v0.2.1
ENV TARBALL=dvi-to-text_${VERSION}_x86_64-unknown-linux-musl.tar.xz
# reitzig/texlive-base does not have curl
RUN wget https://github.com/jared-hughes/dvi-to-text/releases/download/${VERSION}/${TARBALL} && \
  tar -xJf ${TARBALL} -C /

# Add a full 8-bit font to allow emitting non-ascii bytes
COPY octet.mf /usr/local/texlive/2022/texmf-dist/fonts/source/public/cm
# update kpathsea file database
RUN ls -R /usr/local/texlive/2022/texmf-dist > /usr/local/texlive/2022/texmf-dist/ls-R && \
  # remove all ligatures
  rm -r /usr/local/texlive/2022/texmf-dist/fonts/tfm/public/cm/*.tfm && \
  sed -i 's/ligs:=./ligs:=0/g' /usr/local/texlive/2022/texmf-dist/fonts/source/public/cm/*.mf && \
  # load octet font
  echo "\font\octet=octet" >> /usr/local/texlive/2022/texmf-dist/tex/plain/base/plain.tex && \
  # apply changes
  mktexfmt tex.fmt

# Prune unneeded stuff
RUN bash -c "rm -r /usr/local/texlive/2022/texmf-dist/{doc,source,scripts,bibtex,dvips}" && \
  bash -c "rm -r /usr/local/texlive/2022/texmf-dist/{tex,fonts,xdvi,metafont,makeindex,mft,dvipdfmx}" && \
  bash -c "rm -r /usr/local/texlive/2022/bin/x86_64-linuxmusl/{lua*,pdf*,*dvi*,mf-nowin,bibtex}" && \
  bash -c "rm -r /usr/local/texlive/2022/texmf-var/web2c/{lua*,pdftex,metafont}" && \
  rm -r /usr/local/texlive/2022/texmf-var/fonts && \
  rm -r /usr/local/texlive/2022/tlpkg

FROM codegolf/lang-base

COPY --from=builder /bin                      /bin
COPY --from=builder /lib/ld-musl-x86_64.so.1  /lib/ld-musl-x86_64.so.1
COPY --from=builder /dvi-to-text              /usr/bin/
COPY --from=builder /usr/lib/libncursesw.so.6 \
                    /usr/lib/libreadline.so.8 /usr/lib/
COPY --from=builder /usr/local/bin            /usr/local/bin
COPY --from=builder /usr/local/texlive/2022   /usr/local/texlive/2022

COPY tex /usr/bin/

ENTRYPOINT ["/usr/bin/tex"]

CMD ["-v"]
