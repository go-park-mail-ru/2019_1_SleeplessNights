#!/bin/bash

#TODO CREATE

#Запускаем consul

consul_config_dir=$BASEPATH/consul.d

xterm -hold -e "
    consul agent -dev -microservices-dir=$consul_config_dir -node=machine -ui -http-port 8010 -client 0.0.0.0 > /dev/null
" &

for config in $(ls $consul_config_dir)
do
ms_dir=${config%.*} #Обрезаем расширение файла
port =
done
