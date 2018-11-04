# Staffio

An OAuth2 server with management for enterprise employees.


## features:

* All employees in LDAP.
* Login and general member settings.
* Reset password with email.
* Create, Edit and Remove employees with special manager.
* Client ID and Secret of all clients maintenance.
* Simplified content management for aritcles and links.
* A general OAuth2 authentication and authorization provider.
* Directly CAS implement for v1 and V2.

## APIs of oauth2

### Authorize (browse page)
> GET | POST /authorize

### Retrieve Token
> GET | POST /token

### Get Info
> GET | POST /info/{topic}


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


## prepare development

### checkout

````sh
mkdir -p $GOPATH/src/github.com/liut
cd $GOPATH/src/github.com/liut
git clone https://github.com/liut/staffio.git
cd $GOPATH/src/liut/staffio
make dep
````

### LDAP

#### append schema first time only

- Special schema: [ldif](database/ldap_schema/staffio.ldif) or [schema](database/ldap_schema/staffio.schema)

```sh
cat database/ldap_schema/staffio.schema | sudo slapd-config schema write staffio
```

*TODO*

### database
*recommend docker for development*
````sh
docker run -e DB_NAME=staffio -e DB_USER=staffio -e DB_PASS=mypassword -e TZ=Hongkong -p 54322:5432 -d --name staffio-db lcgc/postgresql:9.5.4
cat database/schema.sql | docker exec -i staffio-db psql -Ustaffio
cat database/init.sql | docker exec -i staffio-db psql -Ustaffio

-- example ldif

ldapadd -x -D "cn=admin,dc=example,dc=org" -W -f database/example/init.ldif

-- demo client
echo "INSERT INTO oauth_client VALUES(1, '1234', 'Demo', 'aabbccdd', 'http://localhost:3000/appauth', '{}', now());" | docker exec -i staffio-db psql -Ustaffio staffio

````

### environment

```
    cp -n .env.example .env
```

> `cat .env`
```
STAFFIO_PREFIX=http://localhost:3000
STAFFIO_PASSWORD_SECRET="mypasswordsecret"
STAFFIO_HTTP_LISTEN="localhost:3000"
STAFFIO_LDAP_HOSTS=slapd.hostname
STAFFIO_LDAP_BASE="dc=example,dc=net"
STAFFIO_LDAP_BIND_DN="cn=admin,dc=example,dc=net"
STAFFIO_LDAP_PASS="myadminpassword"
STAFFIO_SESS_NAME="staff_sess"
STAFFIO_SESS_SECRET="very-secret"
STAFFIO_SESS_MAXAGE=86400
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
scp *-linux-amd64-*.tar.xz remote:/path/of/app/bin/
rsync -rpt --delete templates htdocs remote:/path/of/app/
```

## Plan

* Peoples and groups sync with WxWork
* Signin with WxWork
* Notification system
* Export for backup
* Batch import or restore from backup
