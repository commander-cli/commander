#!/bin/bash

# Execute it viá "make integration" from the project root

cd integration/containers/ssh && docker build -t commander-int-ssh-server -f Dockerfile .
docker run -d -p 2222:22 --name=int-ssh-server commander-int-ssh-server