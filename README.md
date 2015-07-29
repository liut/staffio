# staffio
## define

    A oauth2 server provider for enterprise employees


## prepare

````bash
docker run -e DB_NAME=staffio -e DB_USER=staffio -e DB_PASS=mypassword -d --name staffio-db lcgc/postgresql:9.4.4

PGHOST=`docker inspect -f "{{.NetworkSettings.IPAddress}}" staffio-db`
psql -h $PGHOST -Ustaffio -W staffio < database/schema.sql
psql -h $PGHOST -Ustaffio -W staffio < database/init.sql
````


