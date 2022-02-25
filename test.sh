#!/usr/bin/env bash
set -exo pipefail

# Usage: ./test.sh <make target>
# Execute a makefile target inside a Docker environment

MAKE_TARGET="test-coverage-all-codeclimate"
if [[ -n "$1" ]]; then
  MAKE_TARGET="$1"
fi

NO_CLEANUP=0
if [[ "$2" == "--no-cleanup" ]]; then
  NO_CLEANUP=1
fi

cleanup() {
  echo "Starting cleanup"
  container_name="$1"
  docker stop commander-int-ssh-server || true
  docker rm commander-int-ssh-server || true
  docker network rm commander_test
}

start_ssh_conatiner() {
  local ssh_container_name="commander-int-ssh-server"
  if [[ ! $(docker ps | grep $ssh_container_name) ]]; then
      echo "Starting SSH server"
      docker run -d \
        --rm \
        --ip=172.28.0.2 \
        --network commander_test \
        --name "$ssh_container_name" \
        "$ssh_container_name"
  fi
}

if [[ ! "$(docker ps -a | grep -w commander-int-test)" ]]; then
    docker build -t commander-int-test -f integration/containers/test/Dockerfile .
fi

if [[ ! "$(docker ps -a | grep commander-int-ssh-server)" ]]; then
    docker build -t commander-int-ssh-server -f integration/containers/ssh/Dockerfile .
fi

if [[ ! "$(docker network ls | grep commander_test)" ]]; then
    docker network create --driver=bridge --subnet=172.28.0.0/16 commander_test
fi

start_ssh_conatiner

echo "Starting tests"
docker run \
  --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --network commander_test \
  --name commander-integration-go-test \
  --env CC_TEST_REPORTER_ID="${CC_TEST_REPORTER_ID}" \
  commander-int-test \
  make "$MAKE_TARGET"

status_code="$?"

if [[ "$NO_CLEANUP" -eq "0" ]]; then
  cleanup
fi

exit $status_code

