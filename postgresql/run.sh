#!/usr/bin/env bash

go run ${BASEPATH}/postgresql/get_postgres_parameter.go port > port
port=$(tail -n 1 port)
rm port

go run ${BASEPATH}/postgresql/get_postgres_parameter.go host > host
host=$(tail -n 1 host)
rm host

go run ${BASEPATH}/postgresql/get_postgres_parameter.go db_name > db_name
db_name=$(tail -n 1 db_name)
rm db_name

go run ${BASEPATH}/postgresql/get_postgres_parameter.go user > user
user=$(tail -n 1 user)
rm user

go run ${BASEPATH}/postgresql/get_postgres_parameter.go password > password
password=$(tail -n 1 password)
rm password

image="sleepless_nights_postgres"
docker build ${BASEPATH}/postgresql/. -q -t ${image}

container=$(docker run \
    -e POSTGRES_DB=${db_name} \
    -e POSTGRES_USER=${user} \
    -e POSTGRES_PASSWORD=${password} \
    -h ${host} \
    -p ${port}:5432 \
    -d ${image});

echo "Postgres container: ${container}"
