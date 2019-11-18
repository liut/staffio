FROM lcgc/postgresql:9.6.13

ENV DB_NAME=staffio DB_USER=staffio PG_EXTENSIONS=pg_trgm

ADD staffio_*.sql /docker-entrypoint-initdb.d/

