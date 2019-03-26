[![Build Status](https://travis-ci.org/SimonBaeumer/commander.svg?branch=master)](https://travis-ci.org/SimonBaeumer/commander)
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonBaeumer/commander)](https://goreportcard.com/report/github.com/SimonBaeumer/commander)
[![Maintainability](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/maintainability)](https://codeclimate.com/github/SimonBaeumer/commander/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/test_coverage)](https://codeclimate.com/github/SimonBaeumer/commander/test_coverage)
[![Github All Releases](https://img.shields.io/github/downloads/SimonBaeumer/commander/total.svg)](https://github.com/SimonBaeumer/commander/releases)

# Commander

Define `YAML` based test suites for your command line applications.

- [Quick start](#quick-start)
- [Installation](#installation)
  * [Linux & osx](#linux---osx)
  * [Windows](#windows)
- [Example](#example)
- [Minimal test](#minimal-test)
- [Usage](#usage)
- [Development](#development)

## Quick start

 - [Install](#installation) `commander` and add it to your `path` environment variable
 - Create a `commander.yaml` in your project root
 - Add a [minimal test](#minimal-test)
 - Run `./commander test`
 
For more information take a look at the [manual](docs/manual.md), the [examples](examples) or the [integration tests](integration).

## Installation

### Linux & osx

```bash
# Install latest version to /usr/local/bin
curl -fsSL https://raw.githubusercontent.com/SimonBaeumer/commander/master/install.sh | sh

# Install v0.1.0 version to ~/bin
curl -fsSL https://raw.githubusercontent.com/SimonBaeumer/commander/master/install.sh | COMMANDER_VER=v0.1.0 COMMANDER_DST=~/bin sh
```

### Windows

 - Download the current [release](https://github.com/SimonBaeumer/commander/releases/latest)
 - Add the path to your [path](https://docs.alfresco.com/4.2/tasks/fot-addpath.html) environment variable
 - Test it: `commander --version`

## Example

You can find more examples in `examples/`

```
# Build the project
$ make build

# Execute test suite
Starting test file examples/commander.yaml...

✓ it should print hello world
✓ echo hello
✓ it should validate exit code
✓ it should fail

Duration: 0.005s
Count: 4, Failed: 0
```

## Minimal test

```yaml
tests:
    echo hello:
        stdout: hello
        exit-code: 0
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
   --help, -h     show help
   --version, -v  print the version

```


## Development

```
# Initialise dev env
$ make init

# Build the project binary
$ make build

# Unit tests
$ make test

# Coverage
$ make test-coverage

# Integration tests
$ make test-integration

# Add depdencies to vendor
$ make deps
```
