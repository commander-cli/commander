[![Build Status](https://travis-ci.org/SimonBaeumer/commander.svg?branch=master)](https://travis-ci.org/SimonBaeumer/commander)
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonBaeumer/commander)](https://goreportcard.com/report/github.com/SimonBaeumer/commander)

# Commander

Testing tool for command line applications.

## Usage

```
$ make build
$ ./commander ./example/commander.yaml
✓  more printing
✓  it should print hello world
✓  it should print something
```

## Todo:
 - suite fails -> error exit code
 - logging / verbose output
 - print errors in colors
 - execute a single test

 - go api
 - command execution
   - environment variables
   - arguments?
   - timeout
 - exit code *done*
 - stdout
    - Validate against string *done*
    - Validate against file
    - Validate against line
    - Validate with wildcards / regex
    - Validate against template
 - stderr
    - Validate against string *done*
    - Validate against file
    - Validate with wildcards
    - Validate against template
 - testing interactive applications?
 - Support different os
   - Windows
   - MacOs
   - Linux
      - debian
      - ubuntu
      - rhel
      - centos
      - alpine
      
## Open

 - support for...
    - docker
    - docker-compose
    - lxc
    - vagrant

## Architecture

 - runtime?
     - test-executer
     - ordering?
 - interpreter?
