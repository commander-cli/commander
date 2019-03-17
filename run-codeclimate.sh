#!/usr/bin/env bash

if [[ -z $1 ]]; then
    echo "No argument given"
fi

docker pull codeclimate/codeclimate-structure
docker pull codeclimate/codeclimate-duplication

docker run \
  --interactive --tty --rm \
  --env CODECLIMATE_CODE="$PWD" \
  --volume "$PWD":/code \
  --volume /var/run/docker.sock:/var/run/docker.sock \
  --volume /tmp/cc:/tmp/cc \
  codeclimate/codeclimate $@