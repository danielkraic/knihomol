#!/bin/bash

CONTAINER_NAME="knihomoldb"
IS_RUNNING=$(docker container ls | grep $CONTAINER_NAME)

function knihomol_status() {
    if [ -n "$IS_RUNNING" ]; then
        echo "$CONTAINER_NAME container is running"
        exit 0
    fi

    IS_STOPPED=$(docker container ls -a | grep $CONTAINER_NAME)
    if [ -n "$IS_STOPPED" ]; then
        echo "$CONTAINER_NAME container is stopped"
        exit 0
    fi

    echo "$CONTAINER_NAME container not existing"
}

function knihomol_start() {
    if [ -n "$IS_RUNNING" ]; then
        echo "$CONTAINER_NAME container is already running"
        exit 0
    fi

    echo "$CONTAINER_NAME container is NOT running"
    if ! docker run -d -p 27017:27017  -v /opt/knihomoldb:/data/db --name $CONTAINER_NAME mongo; then
        echo "failed to start $CONTAINER_NAME container"
        exit 1
    fi

    echo "$CONTAINER_NAME container was started"
}

function knihomol_stop() {
    if [ -z "$IS_RUNNING" ]; then
        echo "$CONTAINER_NAME container is NOT running"
        exit 0
    fi

    if ! docker container stop $CONTAINER_NAME; then
        echo "failed to start $CONTAINER_NAME"
        exit 1
    fi

    echo "$CONTAINER_NAME container was stopped"
}

function knihomol_rm() {
    if [ -n "$IS_RUNNING" ]; then
        knihomol_stop
    fi

    IS_STOPPED=$(docker container ls -a | grep $CONTAINER_NAME)
    if [ -z "$IS_STOPPED" ]; then
        echo "$CONTAINER_NAME container not existing"
        exit 0
    fi

    if ! docker container rm $CONTAINER_NAME; then
        echo "failed to rm $CONTAINER_NAME container"
        exit 1
    fi

    echo "$CONTAINER_NAME container was removed"
}

case "$1" in
    status)
        knihomol_status
    ;;

    start)
        knihomol_start
    ;;

    stop)
        knihomol_stop
    ;;

    rm)
        knihomol_rm
    ;;

    *)
        echo "Usage: $0 status|start|stop|rm" >&2
        exit 1
    ;;
esac
