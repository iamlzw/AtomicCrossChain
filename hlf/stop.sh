#!/bin/sh
docker-compose -f orderer.yaml -f peer.yaml down
docker volume prune
