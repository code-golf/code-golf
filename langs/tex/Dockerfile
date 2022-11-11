FROM rust:1.65 as builder
WORKDIR /dvi-to-text
RUN USER=root cargo install dvi-to-text --version 0.1.0

# Not sure if texlive should be tagged as well
# (see build-langs `Tag the builder stage so "docker system prune" won't remove it.`)

# TODO: pin version
FROM texlive/texlive:latest-small as texlive

FROM codegolf/lang-base

COPY --from=builder /usr/local/cargo/bin/dvi-to-text /usr/bin/

# TODO: slim. If we want just plain TeX, this could probably be <30 MB

COPY --from=texlive /bin   /bin
COPY --from=texlive /lib64   /lib64
COPY --from=texlive /lib/x86_64-linux-gnu /lib/x86_64-linux-gnu
COPY --from=texlive /usr/local/texlive/2022 /usr/local/texlive/2022
COPY tex /usr/bin/

ENTRYPOINT ["/usr/bin/tex"]

CMD ["-v"]
