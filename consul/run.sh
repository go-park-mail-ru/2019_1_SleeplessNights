#!/bin/bash
config_file=${BASEPATH}/config.json
config=$(cat ${config_file}) #Подгружаем конфиг консула
client_addr=$(jq -r '.client_addr' ${config_file})
http_port=$(jq -r '.ports.http' ${config_file})  #Достаём из конфига порты, которые нужно прокинуть наружу

#TODO start consul server

docker run \
    -d --net=host \
    -e "CONSUL_LOCAL_CONFIG=${config}" \
    consul agent -ui -dev;

go run ${BASEPATH}/update_kv.go -addr ${client_addr}:${http_port}
