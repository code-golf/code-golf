FROM alpine:3.23 AS builder

ENV BIN='/usr/bin/golfscript' VER=6155e9f

RUN wget -qO $BIN https://raw.githubusercontent.com/lynn/golfscript/$VER/golfscript.rb \
 && chmod +x $BIN

# Shebangs can only have one operand, "encoding ASCII" → "encoding=ASCII".
RUN sed -i '1c#!/usr/bin/ruby --encoding=ASCII-8BIT' $BIN

FROM codegolf/lang-ruby

COPY --from=0 /usr/bin/golfscript /usr/bin/

ENTRYPOINT ["golfscript"]
