FROM alpine:3.22 AS builder

ENV VER=6155e9f

RUN wget -O- https://github.com/lynn/golfscript/tarball/$VER \
  | tar xz --strip-components 1

RUN chmod +x golfscript.rb

# Shebangs can only have one operand, "encoding ASCII" → "encoding=ASCII".
RUN sed -i '1c#!/usr/bin/ruby --encoding=ASCII-8BIT' golfscript.rb

FROM codegolf/lang-ruby

COPY --from=0 /golfscript.rb /usr/bin/golfscript

ENTRYPOINT ["golfscript"]
