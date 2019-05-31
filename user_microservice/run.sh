#!/usr/bin/env bash
name="user_ms"
container=$(docker run \
    -e "CONSUL_ADDR=${CONSUL_ADDR}" \
    --name ${name} --rm \
    --network host \
    ${BACKEND_IMAGE} go run /server/user_microservice/main.go);

echo "User-MS container: ${container}"