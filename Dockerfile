FROM alpine:3.4
MAINTAINER Eagle Liut <eagle@dantin.me>

ENV VERSION=v0.2.1 \
    PGHOST="staffio-db" \
    STAFFIO_HTTP_LISTEN=":80" \
    STAFFIO_LDAP_HOST="slapd" \
    STAFFIO_LDAP_BASE="dc=example,dc=org" \
    STAFFIO_ROOT="/app"
ENV DOWNLOAD_URL https://github.com/liut/staffio/releases/download/$VERSION/staffio-linux-amd64-$VERSION.tar.xz

RUN apk add --virtual build-dependencies --update \
  curl \
  ca-certificates \
  && curl -L $DOWNLOAD_URL | tar Jxv -C /usr/bin \
  && apk del build-dependencies \
  && rm -rf /var/cache/apk/*

RUN mkdir /app
WORKDIR /app

ADD templates /app/templates

EXPOSE 80

# ENTRYPOINT ["/usr/bin/staffio"]
CMD ["/usr/bin/staffio"]
