#!/usr/bin/env bash

#current_dir=$(pwd)
#cd ${BASEPATH}/chat_microservice
#go build -tags netgo -a -v .
#cd ${current_dir}

#image="sleepless_nights_chat_microservice"
#docker build ${BASEPATH}/chat_microservice/. -t ${image}

name="chat_ms"
container=$(docker run \
    -e "CONSUL_ADDR=${CONSUL_ADDR}" \
    --name ${name} --rm \
    --network host \
    -d ${BACKEND_IMAGE} go run /server/chat_microservice/main.go);

echo "Chat-MS container: ${container}"
