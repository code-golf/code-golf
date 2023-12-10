FROM clojure:temurin-21-tools-deps-alpine as builder

RUN chown -R nobody: /root

FROM codegolf/lang-base

ENV JAVA_HOME /opt/java/openjdk

COPY --from=0 / /

ENTRYPOINT ["clojure"]

CMD ["-M", "-e", "(println (clojure-version))"]
