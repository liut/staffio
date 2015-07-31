FROM alpine:3.2
MAINTAINER Eagle Liut <eagle@dantin.me>

ENV STAFFIO_LDAP_HOST "slapd"
ENV STAFFIO_LDAP_BASE "dc=example,dc=org"
ENV STAFFIO_HTTP_LISTEN ":80"
ENV STAFFIO_ROOT "/app"
ENV PGHOST "staffio-db"

RUN mkdir /app
WORKDIR /app

ADD dist-tight/linux/amd64/staffio /usr/bin/

ADD templates /app/templates

EXPOSE 80

# ENTRYPOINT ["/usr/bin/staffio"]
CMD /usr/bin/staffio
