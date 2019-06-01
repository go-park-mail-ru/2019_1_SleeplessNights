#!/usr/bin/env bash
#Разбираем параметры командной строки
static="${BASEPATH}/main_microservice/static"
while [[ -n "$1" ]]
do
case "$1" in
--prod) static="/usr/share/nginx/html/assets";;
*) return 1;;
esac
shift
done

name="main_ms"
container=$(docker run \
    -e "CONSUL_ADDR=${CONSUL_ADDR}" \
    --name ${name} --rm \
    --network host \
    -v ${BASEPATH}/main_microservice/static:${static} \
    -d ${BACKEND_IMAGE} go run /server/main_microservice/main.go);

echo "Main-MS container: ${container}"
