#!/bin/sh

# shellcheck disable=SC2046
eval $(minikube docker-env)
docker build -t gof-deploy:1.0 .