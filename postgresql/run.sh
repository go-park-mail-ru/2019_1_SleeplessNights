#!/usr/bin/env bash

go run ${BASEPATH}/consul/helpers/get_parameter.go postgres.port > port
port=$(tail -n 1 port)
rm port

go run ${BASEPATH}/consul/helpers/get_parameter.go postgres.host > host
host=$(tail -n 1 host)
rm host

go run ${BASEPATH}/consul/helpers/get_parameter.go postgres.db_name > db_name
db_name=$(tail -n 1 db_name)
rm db_name

go run ${BASEPATH}/consul/helpers/get_parameter.go postgres.user > user
user=$(tail -n 1 user)
rm user

go run ${BASEPATH}/consul/helpers/get_parameter.go postgres.password > password
password=$(tail -n 1 password)
rm password

name="postgresql"
container=$(docker run \
    --name ${name} --rm \
    -e POSTGRES_DB=${db_name} \
    -e POSTGRES_USER=${user} \
    -e POSTGRES_PASSWORD=${password} \
    -h ${host} \
    -p ${port}:5432 \
    -v ${BASEPATH}/postgresql/data:/var/lib/postgresql/data \
    -v ${BASEPATH}/postgresql/dump.sql:/docker-entrypoint-initdb.d/dump.sql \
    -d postgres);

echo "Postgres container: ${container}"
