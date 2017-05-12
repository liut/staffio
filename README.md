# staffio
## define

    An OAuth2 server with management for enterprise employees


## features

* All employees in LDAP
* Login and general member settings
* Reset password with email
* Create, Edit and Remove employees with special manager
* Client ID and Secret of all clients maintenance
* Simplified content management of aritcles and links
* A general OAuth2 authentication and authorization provider
* Directly CAS implement for v1 and V2


### checkout

````sh
mkdir -p $GOPATH/src/lcgc/platform
cd $GOPATH/src/lcgc/platform
git clone https://github.com/liut/keeper.git
git clone https://github.com/liut/staffio.git
````

## prepare

### LDAP

*TODO*

### database
*recommend docker for development*
````sh
docker run -e DB_NAME=staffio -e DB_USER=staffio -e DB_PASS=mypassword -e TZ=Hongkong -p 54322:5432 -d --name staffio-db lcgc/postgresql:9.5.4
cat database/schema.sql | docker exec -i staffio-db psql -Ustaffio
cat database/init.sql | docker exec -i staffio-db psql -Ustaffio

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

## TODO

* Peoples and groups sync with WxWork
* Signin with WxWork
* Notification system
* Export for backup
* Batch import or restore from backup
