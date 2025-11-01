FROM codegolf/lang-base

COPY --from=codegolf/binutils \
    /usr/bin/aarch64-linux-gnu-as \
    /usr/bin/aarch64-linux-gnu-ld \
    /usr/bin/aarch64-linux-gnu-objcopy \
    /usr/bin/

COPY --from=codegolf/qemu \
    /usr/bin/as-ld-qemu \
    /usr/bin/qemu-aarch64 \
    /usr/bin/

ENTRYPOINT ["qemu-aarch64"]

CMD ["--version"]
