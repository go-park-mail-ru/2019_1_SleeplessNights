#!/usr/bin/env bash
name="game_ms"
container=$(docker run \
    -e "CONSUL_ADDR=${CONSUL_ADDR}" \
    --name ${name} --rm \
    --network host \
    -d ${BACKEND_IMAGE} go run /server/game_microservice/main.go);

echo "Game-MS container: ${container}"
