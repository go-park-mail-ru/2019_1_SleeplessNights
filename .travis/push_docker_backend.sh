#!/usr/bin/env bash

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker build . -t sleepless_nights_backend
docker push uimin1maksim/sleepless_nights_backend
