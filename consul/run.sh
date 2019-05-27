#!/bin/bash
config_file=${BASEPATH}/consul/config.json
config=$(cat ${config_file}) #Подгружаем конфиг консула
client_addr=$(jq -r '.client_addr' ${config_file})
http_port=$(jq -r '.ports.http' ${config_file})  #Достаём из конфига порты, которые нужно прокинуть наружу
https_port=$(jq -r '.ports.https' ${config_file})
grpc_port=$(jq -r '.ports.grpc' ${config_file})

#TODO start consul server

container=$(docker run \
    -d \
    -p ${http_port}:${http_port} \
    -p ${https_port}:${https_port} \
    -p ${grpc_port}:${grpc_port} \
    -h ${client_addr} \
    -e "CONSUL_LOCAL_CONFIG=${config}" \
    consul agent -ui -dev);

echo "Consul container: ${container}"

go run ${BASEPATH}/consul/update_kv.go -addr ${CONSUL_ADDR}