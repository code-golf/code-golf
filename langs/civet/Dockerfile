FROM alpine:3.20 AS builder

RUN apk add --no-cache npm

RUN npm i -g @danielx/civet

ENTRYPOINT ["civet"]

CMD ["-v"]
