#!/usr/bin/env bash
name="main_ms"
container=$(docker run \
    -e "CONSUL_ADDR=${CONSUL_ADDR}" \
    --name ${name} --rm \
    --network host \
    -d ${BACKEND_IMAGE} go run /server/main_microservice/main.go);

echo "Main-MS container: ${container}"
