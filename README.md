[![Build Status](https://travis-ci.org/SimonBaeumer/commander.svg?branch=master)](https://travis-ci.org/SimonBaeumer/commander)
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonBaeumer/commander)](https://goreportcard.com/report/github.com/SimonBaeumer/commander)
[![Maintainability](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/maintainability)](https://codeclimate.com/github/SimonBaeumer/commander/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/test_coverage)](https://codeclimate.com/github/SimonBaeumer/commander/test_coverage)
[![Github All Releases](https://img.shields.io/github/downloads/SimonBaeumer/commander/total.svg)](https://github.com/SimonBaeumer/commander/releases)

# Commander

Define `YAML` based test suites for your command line applications.

## Installation

```bash
# Install latest version to /usr/local/bin
curl -fsSL https://raw.githubusercontent.com/SimonBaeumer/commander/master/install.sh | sh

# Install v0.1.0 version to ~/bin
curl -fsSL https://raw.githubusercontent.com/SimonBaeumer/commander/master/install.sh | COMMANDER_VER=v0.1.0 COMMANDER_DST=~/bin sh
```

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
   Commander - CLI app testing

USAGE:
   commander [global options] command [command options] [arguments...]

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
config:
    env:
    - ENV_KEY=value
    dir: /tmp #set the current working dir
tests:
    it will print hello world:
        cmd: echo hello world
        stdout:
            lines:
                1: hello world
            contains: 
                - hello world
            exactly: hello world
        exit-code: 0
            
    it will print hello:
        cmd: echo hello
        stdout: hello
        exit-code: 0
        
    it prints variable:
        cmd: echo $ENV_KEY
        stdout: value
        exit-code: 0
    
    it overwrites dir:
        cmd: pwd
        config:
            dir: /home/commander
        stdout: /home/commander
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
