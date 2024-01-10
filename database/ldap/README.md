# Build docker image for ldap service


## build a docker image

    make build && make build-tag


## run a docker container

	make run base_dn="dc=mydomain,dc=net" password=mypassword

	or

    docker run --name ldap \
    	-e LDAP_BASE_DN="dc=mydomain,dc=net" \
    	-e LDAP_ADMIN_PASSWORD="mysecret" \
    	-p 1389:389 -p 1636:636 \
    	-d liut7/staffio-ldap:2.4

## write schema into slapd config (local only)

```sh

sudo ./sbin/slapd-config schema write staffio < schema/staffio.schema

```
