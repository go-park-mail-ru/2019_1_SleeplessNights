#!/bin/bash

#Запускаем consul

consul_config_dir=$BASEPATH/consul.d
system_config=$BASEPATH/system-config.json

if [[ -d ${consul_config_dir} ]]
then
    if [[ -r ${system_config} ]]
    then
    node=$(jq -r '.consul.node' ${system_config})
    http_port=$(jq -r '.consul.http_port' ${system_config})
    client=$(jq -r '.consul.client' ${system_config})
    consul agent -dev -config-dir=$consul_config_dir -ui -node=$node -http-port $http_port -client $client
    else
    echo "Can't find BASEPATH/system-microservices.json = $system_config"
    exit 1
    fi
else
echo "Can't find directory BASEPATH/consul.d = $consul_config_dir"
exit 1
fi