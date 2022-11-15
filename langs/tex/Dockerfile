FROM rust:1.65-alpine as builder
WORKDIR /dvi-to-text
RUN USER=root cargo install dvi-to-text --version 0.2.0

# Not sure if texlive should be tagged as well
# (see build-langs `Tag the builder stage so "docker system prune" won't remove it.`)

FROM reitzig/texlive-base:2022.7 as texlive

# Add a full 8-bit font to allow emitting non-ascii bytes
COPY octet.mf /usr/local/texlive/2022/texmf-dist/fonts/source/public/cm
RUN ls -R /usr/local/texlive/2022/texmf-dist > /usr/local/texlive/2022/texmf-dist/ls-R && \
  echo "\font\octet=octet" >> /usr/local/texlive/2022/texmf-dist/tex/plain/base/plain.tex && \
  mktexfmt tex.fmt

# Prune unneeded stuff
RUN bash -c "rm -r /usr/local/texlive/2022/texmf-dist/{doc,source,scripts,bibtex,dvips}" && \
  bash -c "rm -r /usr/local/texlive/2022/texmf-dist/{tex,fonts,xdvi,metafont,makeindex,mft,dvipdfmx}" && \
  bash -c "rm -r /usr/local/texlive/2022/bin/x86_64-linuxmusl/{lua*,pdf*,*dvi*,mf-nowin,bibtex}" && \
  bash -c "rm -r /usr/local/texlive/2022/texmf-var/web2c/{lua*,pdftex,metafont}" && \
  rm -r /usr/local/texlive/2022/texmf-var/fonts && \
  rm -r /usr/local/texlive/2022/tlpkg

FROM codegolf/lang-base

COPY --from=builder /usr/local/cargo/bin/dvi-to-text /usr/bin/
COPY --from=texlive /bin   /bin
COPY --from=texlive /lib/ld-musl-x86_64.so.1 /lib/ld-musl-x86_64.so.1
COPY --from=texlive /usr/local/bin /usr/local/bin
COPY --from=texlive /usr/local/texlive/2022 /usr/local/texlive/2022
COPY --from=texlive /usr/lib/libreadline.so.8 /usr/lib/libreadline.so.8
COPY --from=texlive /usr/lib/libncursesw.so.6 /usr/lib/libncursesw.so.6
COPY tex /usr/bin/

ENTRYPOINT ["/usr/bin/tex"]

CMD ["-v"]
