FROM alpine:3.5
MAINTAINER Eagle Liut <eagle@dantin.me>

RUN apk add --update \
  bash \
  su-exec \
  && rm -rf /var/cache/apk/*

ENV VERSION=v0.8.6 \
    PGHOST="staffio-db" \
    STAFFIO_HTTP_LISTEN=":3030" \
    STAFFIO_LDAP_HOST="slapd" \
    STAFFIO_LDAP_BASE="dc=example,dc=org" \
    STAFFIO_PASSWORD_SECRET=vajanuyogohusopekujabagaliquha \
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
ADD entrypoint.sh /app/entrypoint.sh

EXPOSE 3030

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["web"]
