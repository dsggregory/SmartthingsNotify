FROM golang:1.11-alpine

WORKDIR /usr/local/smartthings_notify

RUN set -ex && \
    apk update && apk upgrade && \
    apk add ca-certificates gcc git make libc-dev bash tzdata && \
    apk add mariadb mariadb-client

COPY WebApp/ ./

EXPOSE 8080

# socket path used by the mariadb that we install
ENV MYSQL_SOCKET=/run/mysqld/mysqld.sock
ENV DB_DATA_PATH=/var/lib/mysql

RUN sh setup.sh && make

ENTRYPOINT ["/bin/bash", "/usr/local/smartthings_notify/run.sh"]
#ENTRYPOINT ["/bin/bash"]
