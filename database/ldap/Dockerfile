FROM alpine:3.18

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
	&& apk add --update \
	bash \
	openldap \
	openldap-backend-all \
	openldap-clients \
	openldap-passwd-pbkdf2 \
	openssl \
  && rm -rf /var/cache/apk/*

ENV LDAP_LOG_LEVEL=256 LDAP_ORGANISATION="Example Inc." LDAP_DOMAIN=example.org \
	LDAP_BASE_DN="dc=example,dc=org" \
	LDAP_ADMIN_NAME=admin LDAP_ADMIN_PASSWORD=mypassword

EXPOSE 389 636

ADD schema/* /etc/openldap/schema/
ADD sbin/* /usr/sbin/

CMD ["/usr/sbin/start.sh"]
