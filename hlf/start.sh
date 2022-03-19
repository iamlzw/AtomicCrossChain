#!/bin/sh
docker-compose -f orderer.yaml -f peer.yaml up -d
