FROM alpine:3.22 AS builder

RUN apk add --no-cache bash coreutils curl gcc musl-dev

ENV VER=3.7.1

RUN curl -#L https://github.com/scala/scala3/releases/download/$VER/scala3-$VER.tar.gz \
  | tar xz --directory /usr --strip-components 1

COPY scala.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/scala -s -static scala.c

FROM codegolf/lang-java

COPY --from=0 /bin/bash /bin/uname             /bin/
COPY --from=0 /usr/bin/dirname                 \
              /usr/bin/env                     \
              /usr/bin/scala                   /usr/bin/
COPY --from=0 /usr/lib/libacl.so.1             \
              /usr/lib/libattr.so.1            \
              /usr/lib/libcrypto.so.3          \
              /usr/lib/libncursesw.so.6        \
              /usr/lib/libreadline.so.8        \
              /usr/lib/libskarnet.so.2.14      \
              /usr/lib/libutmps.so.0.1         /usr/lib/
COPY --from=0 /usr/libexec/cli-common-platform \
              /usr/libexec/common-shared       \
              /usr/libexec/scala-cli.jar       /usr/libexec/

ENV JAVA_HOME='opt/java'
SHELL ["/bin/bash", "-c"]
# Place all downloaded dependencies somewhere, in the runner tell it where to look, revert to scripting, set --server to false, and done!
RUN /opt/java/bin/java -jar /usr/libexec/scala-cli.jar --power bloop && false

# COPY --from=0 /bin/ls /bin/
# RUN find / -name 'bloop' && false
# RUN ls -a /.local/share/scalacli/bloop/daemon/socket && false

# COPY code.scala /
# RUN /opt/java/bin/java -jar /usr/libexec/scala-cli.jar run code.scala -- abc def ghi 123 456 789

ENTRYPOINT ["scala"]
