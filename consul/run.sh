#!/bin/bash
config_file=${BASEPATH}/consul/config.json
config=$(cat ${config_file}) #Подгружаем конфиг консула
client_addr=$(jq -r '.client_addr' ${config_file})
http_port=$(jq -r '.ports.http' ${config_file})  #Достаём из конфига порты, которые нужно прокинуть наружу

#TODO start consul server

container=$(docker run \
    -d --net=host \
    -e "CONSUL_LOCAL_CONFIG=${config}" \
    consul agent -ui -dev);

echo "Consul container: ${container}"

go run ${BASEPATH}/consul/update_kv.go -addr ${client_addr}:${http_port}