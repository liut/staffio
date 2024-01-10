IMAGE = liut7/staffio-ldap

org = "Example Inc."
base_dn = "dc=example,dc=org"
password = "mysecret"

.PHONY: build run

build:
	docker build --no-cache -t $(IMAGE):latest -f Dockerfile .

build-tag:
	$(eval VER=$(shell docker run --rm $(IMAGE):latest apk info -e openldap -v | awk '{print substr($$1,10)}'))
	echo "re tag to $(IMAGE):$(VER)"
	docker tag $(IMAGE):latest $(IMAGE):$(VER)

run:
	docker run -e LDAP_ORGANIZATION=$(org) \
		-e LDAP_BASE_DN="$(base_dn)" -e LDAP_ADMIN_PASSWORD="$(password)" \
		-d -p 389:389 -p 636:636 --name staffio-ldap $(IMAGE):stable

build-arm64:
	docker build --no-cache -t $(IMAGE):arm64 -f Dockerfile.arm64 .

build-arm64-tag:
	$(eval VER=$(shell docker run --rm $(IMAGE):arm64 apk info -e openldap -v | awk '{print substr($$1,10)}'))
	echo "re tag to $(IMAGE):$(VER)"
	docker tag $(IMAGE):arm64 $(IMAGE):arm64-$(VER)
