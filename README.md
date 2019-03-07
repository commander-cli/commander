[![Build Status](https://travis-ci.org/SimonBaeumer/commander.svg?branch=master)](https://travis-ci.org/SimonBaeumer/commander)
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonBaeumer/commander)](https://goreportcard.com/report/github.com/SimonBaeumer/commander)
[![Maintainability](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/maintainability)](https://codeclimate.com/github/SimonBaeumer/commander/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/test_coverage)](https://codeclimate.com/github/SimonBaeumer/commander/test_coverage)

# Commander

Define `YAML` based test suites for your command line applications.

## Example

```
# Build the project
$ make build

# Execute testsuite
$ ./commander test examples/commander.yaml
Starting test file examples/commander.yaml...

✓ it should exit with error code
✓ it should print hello world

Duration: 0.005s
Count: 2, Failed: 0
```

## Usage

```
NAME:
   Commander - A new cli application

USAGE:
   commander [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     test     Execute the test suite
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose      More output for debugging [$COMMANDER_VERBOSE]
   --help, -h     show help
   --version, -v  print the version
```

## Example yaml file

```
tests:
    it will print hello world:
        cmd: echo hello world
        stdout:
            lines:
                1: hello world
            contains: 
            - hello world
        exit-code: 0
            
    it will print hello:
        cmd: echo hello
        stdout: hello
        exit-code: 0
        
    echo hello:
        exit-code: 0
        
    echo skip:
        exit-code: 0
```

## Development

```
# Build the project binary
$ make build

# Unit tests
make test

# Coverage
make test-coverage

# Integration tests
make test-integration

# Add depdencies to vendor
make deps
```

## Todo:
 - Assert line-count
 - Verbose output
 - command execution
   - environment variables
   - arguments?
   - timeout
   - interactive commands

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
 - Support different os
   - Windows
   - MacOs
   - Linux
      - debian
      - ubuntu
      - rhel
      - centos
      - alpine