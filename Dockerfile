FROM ubuntu:18.04

ENV PGSQLVER 10
ENV DEBIAN_FRONTEND 'noninteractive'

RUN echo 'Europe/Moscow' > '/etc/timezone'

RUN apt-get -o Acquire::Check-Valid-Until=false update
RUN apt install -y gcc git wget
RUN apt install -y postgresql-$PGSQLVER

RUN wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz
RUN tar -xvf go1.12.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT /usr/local/go
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /server
COPY . .

RUN cd /server
RUN go get -u
RUN cat main_microservice/database/deploy_config.json > main_microservice/database/config.json
ENV BASEPATH "/server"
ENV PORT 8080
EXPOSE $PORT

USER postgres

RUN /etc/init.d/postgresql start &&\
	psql --echo-all --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
	createdb -O docker docker &&\
	psql -d docker -f main_microservice/database/sql.sql &&\
	/etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGSQLVER/main/pg_hba.conf &&\
	echo "listen_addresses='*'" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "fsync = off" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "synchronous_commit = off" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "shared_buffers = 512MB" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "random_page_cost = 1.0" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "wal_level = minimal" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\
	echo "max_wal_senders = 0" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf

EXPOSE 5432

USER root

CMD service postgresql start && go run main.go

#FROM golang AS build
#Устанавливаем PostgreSQl
#ENV PGVER 11
#RUN apt-get update
#RUN apt-get install -y curl gnupg2
#RUN apt-get install -y wget && \
#    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
#RUN echo "deb http://apt.postgresql.org/pub/repos/apt bionic-pgdg main" > /etc/apt/sources.list.d/PostgreSQL.list
#RUN add-apt-repository "deb http://archive.ubuntu.com/ubuntu $(lsb_release -sc) universe"
#RUN apt-get update
#RUN apt-get install -y postgresql-$PGVER
#RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
#RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
#Копируем данные в контейнер
#RUN mkdir /server
#WORKDIR /server
#COPY . /server/
#RUN cd /server
#Поднимаеь ЬД из дампа
#FROM postgres AS release
#USER postgres
#EXPOSE 5432
#RUN service postgresql start
#RUN pg_lsclusters
#RUN psql
#RUN cat database/sql.sql | psql -d postgres -c
#RUN cat database/deploy_config.json > database/microservices.json
#Запускаем сервер
#USER root
#RUN go get -u
#ENV BASEPATH "/server"
#ENV PORT 8080
#EXPOSE $PORT


