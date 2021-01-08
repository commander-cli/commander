#!/usr/bin/env bash
set -eo pipefail

# Usage: ./test.sh <make target>
# Execute a makefile target inside a Docker environment

MAKE_TARGET="test-coverage-all-codeclimate"
if [[ -n "$1" ]]; then
  MAKE_TARGET="$1"
fi

docker stop commander-int-ssh-server || true
docker rm commander-int-ssh-server || true
docker stop commander-integration-go-test || true
docker rm commander-integration-go-test || true
docker network rm commander_test || true

docker build -t commander-int-ssh-server -f integration/containers/ssh/Dockerfile .
docker network create  --driver=bridge --subnet=172.28.0.0/16 commander_test

docker run -d \
  --rm \
  --ip=172.28.0.2 \
  --network commander_test \
  --name commander-int-ssh-server \
  commander-int-ssh-server

docker build -t commander-int-test -f integration/containers/test/Dockerfile .
docker run \
  --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --network commander_test \
  --name commander-integration-go-test \
  --env CC_TEST_REPORTER_ID="${CC_TEST_REPORTER_ID}" \
  commander-int-test \
  make "$MAKE_TARGET"

docker stop commander-int-ssh-server || true
docker rm commander-int-ssh-server || true
docker network rm commander_test || true
