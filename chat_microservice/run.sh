#!/usr/bin/env bash
name="chat_ms"
container=$(docker run \
    -e "CONSUL_ADDR=${CONSUL_ADDR}" \
    --name ${name} --rm \
    --network host \
    -d ${BACKEND_IMAGE} go run /server/chat_microservice/main.go);

echo "Chat-MS container: ${container}"
