FROM alpine:edge
MAINTAINER Eagle Liut <eagle@dantin.me>

ENV VERSION v0.2.1
ENV DOWNLOAD_URL https://github.com/liut/staffio/releases/download/$VERSION/staffio-linux-amd64-$VERSION.tar.xz

RUN apk add --virtual build-dependencies --update \
  curl \
  ca-certificates \
  && curl -L $DOWNLOAD_URL | tar Jxv -C /usr/bin \
  && apk del build-dependencies \
  && rm -rf /var/cache/apk/*

ENV STAFFIO_LDAP_HOST "slapd"
ENV STAFFIO_LDAP_BASE "dc=example,dc=org"
ENV STAFFIO_HTTP_LISTEN ":80"
ENV STAFFIO_ROOT "/app"
ENV PGHOST "staffio-db"

RUN mkdir /app
WORKDIR /app

ADD templates /app/templates

EXPOSE 80

# ENTRYPOINT ["/usr/bin/staffio"]
CMD ["/usr/bin/staffio"]
