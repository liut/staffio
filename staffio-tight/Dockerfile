FROM alpine:3.4
MAINTAINER Eagle Liut <eagle@dantin.me>

ENV PGHOST="staffio-db" \
    STAFFIO_HTTP_LISTEN=":80" \
    STAFFIO_LDAP_HOST="slapd" \
    STAFFIO_LDAP_BASE="dc=example,dc=org" \
    STAFFIO_ROOT="/app"

RUN apk --update add \
  bash \
  ca-certificates \
  postgresql-client \
  tzdata \
  && cp /usr/share/zoneinfo/Hongkong /etc/localtime \
  && echo "Hongkong" >  /etc/timezone \
  && apk del tzdata \
  && rm -rf /var/cache/apk/*

RUN mkdir /app
WORKDIR /app

ADD dist-tight/linux_amd64/staffio /usr/bin/

ADD templates /app/templates

EXPOSE 80

# ENTRYPOINT ["/usr/bin/staffio"]
CMD ["/usr/bin/staffio"]
