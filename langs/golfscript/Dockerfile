FROM alpine:3.16 as builder

RUN apk add --no-cache build-base curl openssl-dev

ENV VERSION=2206dbe

RUN curl -L https://github.com/lynn/golfscript/tarball/$VERSION | tar xz \
 && cd /lynn-golfscript-$VERSION && mv golfscript.rb /

FROM codegolf/lang-ruby

COPY --from=0 /golfscript.rb           /golfscript.rb

ENTRYPOINT ["ruby", "--encoding", "ASCII-8BIT", "/golfscript.rb"]

## CMD ["-e", "puts RUBY_VERSION"]
