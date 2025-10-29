FROM node:25.0.0-alpine3.22

COPY --from=codegolf/binutils-wasm     /wasm/ /wasm/
COPY --from=codegolf/aspp-wasm         /wasm/ /wasm/