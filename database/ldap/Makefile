IMAGE = liut7/staffio-ldap
VERSION = 2.4.48-r1

org = "Example Inc."
base_dn = "dc=example,dc=org"
password = "mysecret"

.PHONY: build run

build:
	docker build -t $(IMAGE):$(VERSION) .
	docker tag  $(IMAGE):$(VERSION) $(IMAGE):2.4

run:
	docker run -e LDAP_ORGANIZATION=$(org) \
		-e LDAP_BASE_DN="$(base_dn)" -e LDAP_ADMIN_PASSWORD="$(password)" \
		-d -p 389:389 -p 636:636 --name staffio-ldap $(IMAGE):$(VERSION)
