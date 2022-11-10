# TODO: pin
FROM texlive/texlive:latest-small as builder

FROM codegolf/lang-base

# TODO: slim. If we want just plain TeX, this could probably be <30 MB

COPY --from=0 /bin   /bin
COPY --from=0 /lib64   /lib64
COPY --from=0 /lib/x86_64-linux-gnu /lib/x86_64-linux-gnu
COPY --from=0 /usr/local/texlive/2022 /usr/local/texlive/2022
COPY tex /usr/bin/

ENTRYPOINT ["/usr/bin/tex"]

CMD ["-v"]
