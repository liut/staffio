# Build docker image for ldap service


## build a docker image

    make build


## run a docker container

	make run base_dn="dc=mydomain,dc=net" password=mypassword

	or

    docker run -e LDAP_ORGANIZATION="LCGC Inc." -e LDAP_BASE_DN="dc=mydomain,dc=net" -e LDAP_ADMIN_PASSWORD="mysecret" -d -p 389:389 -p 636:636 --name ldap liut7/staffio-ldap:2.4.44
