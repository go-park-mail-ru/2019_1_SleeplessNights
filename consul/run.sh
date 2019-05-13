#!/bin/bash
config_file=config.json
config=$(cat ${config_file}) #Подгружаем конфиг консула
client_addr=$(jq -r '.client_addr' ${config_file})
http_port=$(jq -r '.ports.http' ${config_file})      #Достаём из конфига порты, которые нужно прокинуть наружу
https_port=$(jq -r '.ports.https' ${config_file})
grpc_port=$(jq -r '.ports.grpc' ${config_file})

#bind_addr=$(jq -r '.bind_addr' ${config_file})
#docker run \
#    -d --net=host \
#    -e 'CONSUL_LOCAL_CONFIG={"skip_leave_on_interrupt": true}' \
#    consul agent -server -bootstrap_expect=1 -bind=${bind_addr} -node=server-agent -enable_script_checks=true ;

docker run \
    -d --net=host \
    -e "CONSUL_LOCAL_CONFIG=${config}" \
    consul agent -ui -dev;

go run update_kv.go -addr ${client_addr}:${http_port}
