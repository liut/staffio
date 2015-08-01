# staffio
## define

    An OAuth2 server with management for enterprise employees


## prepare

````bash
docker run -e DB_NAME=staffio -e DB_USER=staffio -e DB_PASS=mypassword -d --name staffio-db lcgc/postgresql:9.4.4

PGHOST=`docker inspect -f "{{.NetworkSettings.IPAddress}}" staffio-db`
psql -h $PGHOST -Ustaffio -W staffio < database/schema.sql
psql -h $PGHOST -Ustaffio -W staffio < database/init.sql

-- demo client
psql -h $PGHOST -Ustaffio -W -c "INSERT INTO oauth_client VALUES(1, '1234', 'Demo', 'aabbccdd', 'http://localhost:3000/appauth', '{}', now());" staffio


````


