# staffio
## define

    An OAuth2 server with management for enterprise employees


## Features

* All employees in LDAP
* Login and general member settings
* Reset password with email
* Create, Edit and Remove employees with special manager
* Client ID and Secret of all clients maintenance
* Simplified content management of aritcles and links
* A general OAuth2 authentication and authorization Provider
* Directly CAS implement for v1 and V2


## prepare

### database
````sh
docker run -e DB_NAME=staffio -e DB_USER=staffio -e DB_PASS=mypassword -e TZ=Hongkong -p 54322:5432 -d --name staffio-db lcgc/postgresql:9.5.4
cat database/schema.sql | docker exec -i staffio-db psql -Ustaffio
cat database/init.sql | docker exec -i staffio-db psql -Ustaffio

-- demo client
echo "INSERT INTO oauth_client VALUES(1, '1234', 'Demo', 'aabbccdd', 'http://localhost:3000/appauth', '{}', now());" | docker exec -i staffio-db psql -Ustaffio staffio


````

### checkout

````sh
mkdir -p $GOPATH/src/lcgc/platform
cd $GOPATH/src/lcgc/platform
git clone https://github.com/liut/keeper.git
git clone https://github.com/liut/staffio.git
````


## launch development

````sh
go get -u github.com/ddollar/forego
go get -u github.com/ddollar/rerun

forego start
````
