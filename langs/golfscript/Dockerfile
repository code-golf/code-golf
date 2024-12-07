FROM alpine:3.21 AS builder

RUN apk add --no-cache curl

ENV VERSION=6155e9f

RUN curl -L https://github.com/lynn/golfscript/tarball/$VERSION | tar xz \
 && cd /lynn-golfscript-$VERSION && mv golfscript.rb /

# Make executable.
RUN chmod +x /golfscript.rb

# Shebangs can only have one operand, "encoding ASCII" → "encoding=ASCII".
RUN sed -i '1c#!/usr/bin/ruby --encoding=ASCII-8BIT' /golfscript.rb

FROM codegolf/lang-ruby

COPY --from=0 /golfscript.rb /usr/bin/golfscript

ENTRYPOINT ["golfscript"]
