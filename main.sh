#!/bin/bash
#Разбираем параметры командной строки
rebuild=false
tidy=false
run_ms=true
prod=false
while [[ -n "$1" ]]
do
case "$1" in
--prod) prod=true;;
--rebuild) rebuild=true;;
--tidy) tidy=true;;
--no-ms) run_ms=false;;
*) return 1;;
esac
shift
done

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

go_mod_owner=$(stat -c %U ${BASEPATH}/go.mod)
go_sum_owner=$(stat -c %U ${BASEPATH}/go.sum)
#Оптимизируем зависимости
if ${tidy}
then
    current_dir=$PWD
    cd ${BASEPATH}
    go mod tidy
    cd ${current_dir}
fi

#Билдим бекэнд
BACKEND_IMAGE="uimin1maksim/sleepless_nights_backend"
export BACKEND_IMAGE=${BACKEND_IMAGE}
if ${rebuild}
then
    echo "Building backend container..."
    docker build ${BASEPATH}/. -t ${BACKEND_IMAGE}
    echo "Backend container built!"
fi

#Запускаем микросервисы
if ${run_ms}
then
    docker pull uimin1maksim/sleepless_nights_backend
    #Запускаем User-MS
    ${BASEPATH}/user_microservice/run.sh
    #Запускаем Main-MS
    if ${prod}
    then
        ${BASEPATH}/main_microservice/run.sh --prod
    else
        ${BASEPATH}/main_microservice/run.sh
    fi
    #Запускаем Game-MS
    ${BASEPATH}/game_microservice/run.sh
    #Запускаем Chat-MS
    ${BASEPATH}/chat_microservice/run.sh
fi

#Фикисим изменённые файлы
chown ${go_mod_owner} ${BASEPATH}/go.mod
chown ${go_sum_owner} ${BASEPATH}/go.sum
