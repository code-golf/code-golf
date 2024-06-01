FROM rust:1.78-bookworm as builder

RUN cargo install uiua --version 0.10.3

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libgcc_s.so.1 /lib/x86_64-linux-gnu/libgcc_s.so.1
COPY --from=0 /lib/x86_64-linux-gnu/libm.so.6 /lib/x86_64-linux-gnu/libm.so.6
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so.6
COPY --from=0 /lib64/ld-linux-x86-64.so.2 /lib64/ld-linux-x86-64.so.2

COPY --from=0 /usr/local/cargo/bin/uiua /usr/bin/uiua

ENTRYPOINT ["uiua"]

CMD ["--version"]
