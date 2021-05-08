FROM node:16.1.0-alpine3.13 as builder

ENV VERSION=1.0.1

RUN npm install -g @defasm/core@$VERSION

FROM scratch

COPY --from=0 / /
