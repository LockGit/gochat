#!/usr/bin/env bash
# Author : Lock
# Email: lockexit@gmail.com

CONFIG_ENV=$1
IPADDR=$2

valid_ip () {
    local  ip=$1
    local  stat=1

    if [[ $ip =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]]; then
        OIFS=$IFS
        IFS='.'
        ip=($ip)
        IFS=$OIFS
        [[ ${ip[0]} -le 255 && ${ip[1]} -le 255 \
        && ${ip[2]} -le 255 && ${ip[3]} -le 255 ]]
        stat=$?
    fi
    return $stat
}

if [ "${CONFIG_ENV}" = "dev" -o "${CONFIG_ENV}" = "prod" ]
then
    echo "run ${CONFIG_ENV} env..."
else
    echo 'Sorry, error! Please execute For example：'
    echo '    sh run.sh dev 127.0.0.1'
    echo '    sh run.sh prod x.x.x.x'
    exit 1
fi

if ! valid_ip "${IPADDR}";then
    echo "ip address error !!! please check '${IPADDR}' is ip address? "
    exit 1
fi

echo "now env is:${CONFIG_ENV}"
if [ "${CONFIG_ENV}" == "dev" ]
then
    DOCKER_IMAGE="lockgit/gochat:1.18"
elif [ "${CONFIG_ENV}" == "prod" ]
then
    # You can use your own build docker image
    DOCKER_IMAGE="lockgit/gochat:1.18"
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
    -e HOST_IP=${IPADDR} \
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
    -e HOST_IP=${IPADDR} \
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


