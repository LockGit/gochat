#!/usr/bin/env bash
# Author : Lock
# Email: lockexit@gmail.com

echo 'start...'

CONFIG_ENV=$1

if [ ! -n "${CONFIG_ENV}" ]
then
    echo 'Take Me To Your Heart ! Run all environment in a docker container. ^_^'
    echo 'example：'
    echo '    sh run.sh dev'
    echo '    sh run.sh prod'
else
    echo "now env is:${CONFIG_ENV}"
    if [ "${CONFIG_ENV}" == "dev" ]
    then
        DOCKER_IMAGE="lockgit/gochat"
    elif [ "${CONFIG_ENV}" == "prod" ]
    then
        # You can use your own build docker image
        DOCKER_IMAGE="lockgit/gochat"
    fi

    echo "nice, now set docker image is:${DOCKER_IMAGE}"

    # docker run
    if [ "${CONFIG_ENV}" == "dev" ]
    then
        sudo docker run \
        --name gochat-${CONFIG_ENV} \
        -h gochat-${CONFIG_ENV} \
        -e TZ=Asia/Shanghai \
        -e RUN_MODE=dev \
        -e ETCD_DATA_DIR=/root/default.etcd \
        -v `pwd`:/go/src/gochat \
        -v `pwd`/docker/${CONFIG_ENV}/supervisord.d:/etc/supervisord.d \
        -v `pwd`/docker/${CONFIG_ENV}/supervisord.conf:/etc/supervisord.conf \
        -v `pwd`/docker/${CONFIG_ENV}/redis.conf:/root/redis.conf \
        -d \
        -p 8080:8080 \
        -p 7070:7070 \
        -p 7000:7000 \
        -p 7001:7001 \
        -p 7002:7002 \
        ${DOCKER_IMAGE} \
        supervisord -n && docker exec gochat-${CONFIG_ENV} /bin/sh './reload.sh'
    else
        sudo docker run \
        --name gochat-${CONFIG_ENV} \
        -h gochat-${CONFIG_ENV} \
        -e TZ=Asia/Shanghai \
        -e RUN_MODE=prod \
        -e ETCD_DATA_DIR=/root/default.etcd \
        -v `pwd`:/go/src/gochat \
        -v `pwd`/docker/${CONFIG_ENV}/supervisord.d:/etc/supervisord.d \
        -v `pwd`/docker/${CONFIG_ENV}/supervisord.conf:/etc/supervisord.conf \
        -v `pwd`/docker/${CONFIG_ENV}/redis.conf:/root/redis.conf \
        -d \
        -p 8080:8080 \
        -p 7070:7070 \
        -p 7000:7000 \
        -p 7001:7001 \
        -p 7002:7002 \
        ${DOCKER_IMAGE} \
        supervisord -n && docker exec gochat-${CONFIG_ENV} /bin/sh './reload.sh'
    fi

fi
