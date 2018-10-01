FROM alpine:3.6

RUN apk add --update bash openldap openldap-clients openldap-passwd-pbkdf2 \
  && rm -rf /var/cache/apk/*

ENV LDAP_LOG_LEVEL=256 LDAP_ORGANISATION="Example Inc." LDAP_DOMAIN=example.org \
	LDAP_BASE_DN="dc=example,dc=org" \
	LDAP_ADMIN_NAME=admin LDAP_ADMIN_PASSWORD=mypassword

EXPOSE 389 636

ADD schema/staffio.schema /etc/openldap/schema/
ADD scripts/* /usr/sbin/

CMD ["/usr/sbin/start.sh"]
