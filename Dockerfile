FROM alpine

ENV POSTGRES_HOST=localhost \
    POSTGRES_PORT=5432 \
    POSTGRES_USER=checker \
    POSTGRES_PASSWORD=checker \
    POSTGRES_DBNAME=checker

COPY . /usr/local/bin

RUN chmod +x /usr/local/bin/*

EXPOSE 3000

ENTRYPOINT ["docker-entrypoint.sh"]
