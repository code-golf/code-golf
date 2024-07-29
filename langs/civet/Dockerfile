FROM alpine:3.20 AS builder

RUN apk add --no-cache bash curl npm

WORKDIR /usr/src

RUN curl -L https://github.com/DanielXMoore/Civet/archive/2677d5d.tar.gz | tar xz --strip=1

RUN npm i && npm run build

FROM codegolf/lang-base

# next line's beyond diamond but ye, in progress (js might help)
COPY --from=0 . .

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/src                 /usr/local
COPY --from=0 /usr/src/dist            /usr/local/bin

ENTRYPOINT ["civet"]

CMD ["-v"]
