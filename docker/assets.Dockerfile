FROM node:25.0.0-alpine3.22

COPY node_modules ./node_modules
COPY fonts        ./fonts
COPY css          ./css
COPY js           ./js
COPY svg          ./svg
COPY --from=codegolf/binutils-wasm     /wasm/ /wasm/
COPY --from=codegolf/aspp-wasm         /wasm/ /wasm/