FROM alpine:3.2
MAINTAINER Eagle Liut <eagle@dantin.me>

ENV VERSION v0.0.1
ENV DOWNLOAD_URL https://github.com/liut/staffio/releases/download/$VERSION/staffio-linux-amd64-$VERSION.tar.gz

RUN apk add --virtual build-dependencies --update \
  curl \
  ca-certificates \
  && curl -L $DOWNLOAD_URL | tar xvz -C /usr/local/bin \
  && apk del build-dependencies \
  && rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/local/bin/staffio"]
