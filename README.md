# Staffio

An OAuth2 server that provides managed services for enterprise employees.


## Features:

* Employee security information is stored in LDAP.
* Login authentication service and general membership settings.
* Reset password using email and mobile phone number.
* Create edit and delete employees by special members.
* Maintainable APP Client ID and Secret.
* Simple content article and link management.
* Generic OAuth2 authentication and authorization management.
* Directly CAS implement for V1 and V2.


## Objects

### Staff
- `uid`: Username, required
- `cn`: Full Name
- `gn`: FirstName
- `sn`: LastName, required
- `nickname`
- `birthday`: YYYYmmdd
- `gender`: f, m
- `email`: Email
- `mobile`: Cell phone number
- `avatarPath`: Avatar URI
- `description`:
- `joinDate`: YYYYmmdd

### Group
- `name`:
- `description`:
- `members`: []uid

### User (online)
- uid: Username
- name: DisplayName

## APIs of oauth2

### Authorize (browse page)
> GET | POST /authorize

### Retrieve Token
> GET | POST /token

### Get Info
> GET | POST /info/{topic}

#### Info topic
1. `me`: `{me: User}`
2. `me+{groupName}`: `{me: User, group}`
3. `grafana` or `generic`: `{struct for grafana}`

### APIs of <abbr title="Central Authentication Service">CAS</abbr>

| URI | Description |
| -------- | -------- |
| `/login` | credential requestor / acceptor |
| `/logout` | destroy CAS session (logout) |
| `/validate` | service ticket validation |
| `/serviceValidate` | service ticket validation [CAS 2.0] |
| `/proxyValidate` **TODO** | service/proxy ticket validation [CAS 2.0] |
| `/proxy` **TODO** | proxy ticket service [CAS 2.0] |
| `/p3/serviceValidate` **TODO** | service ticket validation [CAS 3.0] |
| `/p3/proxyValidate` **TODO** | service/proxy ticket validation [CAS 3.0] |


## Quick start

### Run all components as docker containers

````sh

# openldap
docker run --name staffio-ldap -p 389:389 -p 636:636 \
	-e LDAP_ADMIN_PASSWORD=mypassword \
	-d liut7/staffio-ldap:latest

# postgresql
docker create --name staffio-db-data -v /var/lib/postgresql busybox:1 echo staffio db data
docker run --name staffio-db -p 54322:5432 \
	-e DB_PASS=mypassword \
	-e TZ=Hongkong \
	--volumes-from=staffio-db-data \
	-d liut7/staffio-db:latest

# staffio main server
docker run --name staffio -p 3030:3030 \
	-e STAFFIO_BACKEND_DSN='postgres://staffio:mypassword@staffio-db/staffio?sslmode=disable' \
	-e STAFFIO_LDAP_HOSTS='ldap://slapd' \
	-e STAFFIO_LDAP_BASE="dc=example,dc=org" \
	-e STAFFIO_LDAP_BIND_DN="cn=admin,dc=example,dc=org" \
	-e STAFFIO_LDAP_PASS='mypassword' \
	--link staffio-db --link staffio-ldap:slapd \
	-d liut7/staffio:latest web

# create a user as first staff and adminstrator
docker exec staffio staffio addstaff -u eagle -p mysecret -n eagleliut --sn liut
docker exec staffio staffio group -g keeper -a eagle

# now can open http://localhost:3030/ in browser

# add a oauth2 client (optional)
docker exec staffio staffio client --add demo --uri http://localhost:3000

# list clients
docker exec staffio staffio client --list

## for testing database
echo "CREATE DATABASE staffiotest WITH OWNER = staffio ENCODING = 'UTF8';" | docker exec -i staffio-db psql -Upostgres
echo "GRANT ALL PRIVILEGES ON DATABASE staffiotest to staffio;" | docker exec -i staffio-db psql -Upostgres

````


## prepare development

### checkout

````sh

go get -u github.com/liut/staffio
cp -n .env.example .env

````

### environment

> `cat .env.example`
```
STAFFIO_HTTP_LISTEN=":3000"
STAFFIO_LDAP_HOSTS=slapd.hostname
STAFFIO_LDAP_BASE="dc=example,dc=org"
STAFFIO_LDAP_BIND_DN="cn=admin,dc=example,dc=org"
STAFFIO_LDAP_PASS="mypassword"
STAFFIO_BACKEND_DSN="postgres://staffio:mypassword@localhost:54322/staffio?sslmode=disable"
STAFFIO_PASSWORD_SECRET="mypasswordsecret"
```

## launch development

````sh
go get -u github.com/ddollar/forego
go get -u github.com/liut/rerun
npm install

forego start
````

## deployment

```sh
make dist package
scp dist/linux_amd64/staffio remote:/opt/staffio/bin/
make fe-build
rsync -rpt --delete templates htdocs remote:/opt/staffio/
```

### add staff
```sh
forego run ./staffio addstaff -u eric -p AF1984 -n George --sn Blair
```

## Plan

* <del>Peoples and groups sync with WxWork</del>
* <del>Signin with WxWork</del>
* Notification
* Export for backup
* Batch import or restore from backup
* I18n
