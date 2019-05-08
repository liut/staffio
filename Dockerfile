FROM golang:1.12
MAINTAINER Eagle Liut <eagle@dantin.me>

ENV GO111MODULE=on GOPROXY=https://goproxy.io
WORKDIR /go/src/github.com/liut/staffio/
COPY main.go go.* ./
COPY pkg ./pkg

RUN pwd && ls \
  && go mod download \
  && CGO_ENABLED=0 GOOS=linux go build -v .


FROM alpine:3.6

RUN apk add --update \
  bash \
  su-exec \
  && rm -rf /var/cache/apk/*

ENV PGHOST="staffio-db" \
    STAFFIO_BACKEND_DSN='postgres://staffio:mypassword@staffio-db/staffio?sslmode=disable' \
    STAFFIO_HTTP_LISTEN=":3030" \
    STAFFIO_LDAP_HOSTS="slapd" \
    STAFFIO_LDAP_BASE="dc=example,dc=org" \
    STAFFIO_LDAP_BIND_DN="cn=admin,dc=example,dc=org" \
    STAFFIO_LDAP_PASS='mypassword' \
    STAFFIO_PASSWORD_SECRET=vajanuyogohusopekujabagaliquha \
    STAFFIO_ROOT="/app"

WORKDIR /app

COPY --from=0 /go/src/github.com/liut/staffio/staffio .
ADD templates /app/templates
ADD entrypoint.sh /app/entrypoint.sh

EXPOSE 3030

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["web"]
