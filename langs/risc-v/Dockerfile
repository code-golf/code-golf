FROM codegolf/lang-base

COPY --from=codegolf/binutils \
    /usr/bin/riscv64-linux-gnu-as \
    /usr/bin/riscv64-linux-gnu-ld \
    /usr/bin/riscv64-linux-gnu-objcopy \
    /usr/bin/

COPY --from=codegolf/qemu \
    /usr/bin/as-ld-qemu \
    /usr/bin/qemu-riscv32 \
    /usr/bin/qemu-riscv64 \
    /usr/bin/

ENTRYPOINT ["qemu-riscv64"]

CMD ["--version"]
