FROM golang:1.18-alpine
MAINTAINER Eagle Liut <eagle@dantin.me>

ENV GO111MODULE=on GOPROXY=https://goproxy.io ROOF=github.com/liut/staffio
WORKDIR /go/src/$ROOF/
COPY main.go go.* ./
COPY htdocs ./htdocs
COPY pkg ./pkg

RUN go mod download \
  && export LDFLAGS="-X ${ROOF}/pkg/settings.buildVersion=$(date '+%Y%m%d')" \
  && env \
  && CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS} -s -w" . \
  && echo "build done"


FROM alpine:3.6

RUN apk add --update \
  bash \
  su-exec \
  && rm -rf /var/cache/apk/*

ENV PGHOST="staffio-db" \
    STAFFIO_BACKEND_DSN='postgres://staffio:mypassword@staffio-db/staffio?sslmode=disable' \
    STAFFIO_HTTP_LISTEN=":3030" \
    STAFFIO_LDAP_HOSTS="ldap://slapd" \
    STAFFIO_LDAP_BASE="dc=example,dc=org" \
    STAFFIO_LDAP_BIND_DN="cn=admin,dc=example,dc=org" \
    STAFFIO_LDAP_PASS='mypassword' \
    STAFFIO_PASSWORD_SECRET=vajanuyogohusopekujabagaliquha \
    STAFFIO_ROOT="/app"

WORKDIR /app

COPY --from=0 /go/src/github.com/liut/staffio/staffio /usr/bin/
ADD templates /app/templates
ADD entrypoint.sh /start.sh

EXPOSE 3030

ENTRYPOINT ["/start.sh"]
CMD ["web"]
