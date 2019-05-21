#!/bin/bash

#Устанавливаем BASEPATH
ABSOLUTE_FILENAME=`readlink -e "$0"` #Абсолютный путь до скрипта
export BASEPATH=`dirname "$ABSOLUTE_FILENAME"`

#Устанавливаем CONSUL_ADDR
consul_config_file=${BASEPATH}/consul/config.json
client_addr=$(jq -r '.client_addr' ${consul_config_file})
http_port=$(jq -r '.ports.http' ${consul_config_file})
export CONSUL_ADDR=${client_addr}:${http_port}

#Запускаем клиент consul
${BASEPATH}/consul/run.sh
#Запускаем postgres
${BASEPATH}/postgresql/run.sh
#Запускаем user_microservice
go run ${BASEPATH}/user_microservice/main.go
#Запускаем main_microservice
#go run ${BASEPATH}/main_microservice/main.go
