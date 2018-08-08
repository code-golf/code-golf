FROM alpine:edge

RUN apk add --no-cache bash binutils clang-dev curl fts-dev gcc libedit-dev make musl-dev

RUN curl -L https://github.com/jsoftware/jsource/archive/90adff2.tar.gz | tar xzf -

RUN mkdir -p ~/git ~/jbld/j64/bin                                  \
 && mv jsource-* ~/git/jsource                                     \
 && cd ~/git/jsource                                               \
 && echo '#define jversion "807"'                > jsrc/jversion.h \
 && echo '#define jplatform "linux"'            >> jsrc/jversion.h \
 && echo '#define jtype "build"'                >> jsrc/jversion.h \
 && echo '#define jlicense "GPL3"'              >> jsrc/jversion.h \
 && echo '#define jbuilder "www.jsoftware.com"' >> jsrc/jversion.h \
 && . make/jvars.sh                                                \
 && make/build_jconsole.sh j64                                     \
 && make/build_libj.sh j64

RUN apk del --no-cache bash binutils clang-dev curl fts-dev gcc make musl-dev

RUN mkdir /j                       \
 && mv ~/git/jsource/jlibrary/* /j \
 && mv /root/jbld/j64/bin/* /j/bin

RUN echo -e "#!/bin/sh -e\n\
\n\
/bin/cat - > /tmp/code.ijs\n\
\n\
exec /j/bin/jconsole /tmp/code.ijs \"\$@\"" > /usr/bin/j && chmod +x /usr/bin/j

ENTRYPOINT ["echo", "807"]
