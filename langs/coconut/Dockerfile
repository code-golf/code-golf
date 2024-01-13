FROM codegolf/lang-base-no-proc

RUN apk add --no-cache build-base curl libffi-dev zlib-dev python3 python3-pip

ENV VER 3.0.4

RUN pip install coconut=$VER

ENTRYPOINT ["coconut"]

CMD ["-v"]
