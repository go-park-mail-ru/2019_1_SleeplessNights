#!/bin/bash
config_file=config.json
config=$(cat ${config_file}) #Подгружаем конфиг консула
http_port=$(jq -r '.ports.http' ${config_file})      #Достаём из конфига порты, которые нужно прокинуть наружу
https_port=$(jq -r '.ports.https' ${config_file})
grpc_port=$(jq -r '.ports.grpc' ${config_file})

docker run \
    -d --net=host \
    --expose ${http_port} --expose ${https_port} --expose ${grpc_port} \
    -e "CONSUL_LOCAL_CONFIG=${config}" \
    consul agent -ui -dev;
