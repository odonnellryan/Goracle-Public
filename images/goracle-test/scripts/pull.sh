#!/bin/bash
rm /var/run/docker.pid >> /dev/null
wrapdocker > /var/log/docker.log 2>&1 &
sleep 5
docker pull docker-test-image