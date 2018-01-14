FROM golang:1.9.2

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8

RUN apt-get update -q && apt-get install -y zip ruby ruby-dev rpm locales && \
  go get github.com/kardianos/govendor && \
  go get github.com/buildkite/github-release && \
  gem install --no-ri --no-rdoc rake fpm package_cloud && \
  echo "en_US UTF-8" > /etc/locale.gen && \
  locale-gen en_US.UTF-8

WORKDIR /go/src/github.com/buildkite/terminal
ADD . /go/src/github.com/buildkite/terminal

CMD [ "make", "dist"]
