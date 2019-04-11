[![Build Status](https://travis-ci.org/SimonBaeumer/commander.svg?branch=master)](https://travis-ci.org/SimonBaeumer/commander)
[![GoDoc](https://godoc.org/github.com/SimonBaeumer/commander?status.svg)](https://godoc.org/github.com/SimonBaeumer/commander)
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonBaeumer/commander)](https://goreportcard.com/report/github.com/SimonBaeumer/commander)
[![Maintainability](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/maintainability)](https://codeclimate.com/github/SimonBaeumer/commander/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/test_coverage)](https://codeclimate.com/github/SimonBaeumer/commander/test_coverage)
[![Github All Releases](https://img.shields.io/github/downloads/SimonBaeumer/commander/total.svg)](https://github.com/SimonBaeumer/commander/releases)

# Commander


Define language independent tests for your command line scripts and programs in simple `yaml` files.

 - It runs on `windows`, `osx` and `linux` 
 - It is a self-contained binary - no need to install a heavy lib or language
 - It is easy and fast to write

For more information take a look at the [manual](docs/manual.md), the [examples](examples) or the [integration tests](integration).

## Installation

### Linux & osx

Visit the [release](https://github.com/SimonBaeumer/commander/releases) page to get the binary for you system. 

```bash
curl -L https://github.com/SimonBaeumer/commander/releases/download/v1.0.0/commander-linux-amd64 -o commander
chmod +x commander
```

### Windows

 - Download the current [release](https://github.com/SimonBaeumer/commander/releases/latest)
 - Add the path to your [path](https://docs.alfresco.com/4.2/tasks/fot-addpath.html) environment variable
 - Test it: `commander --version`


## Quick start

`Commander` will always search for a default `commander.yaml` in the current working directory and execute all defined tests in it.


```bash
# You can even let commander add tests for you!
$ ./commander test examples/commander.yaml
tests:
  echo hello:
    exit-code: 0
    stdout: hello

written to /tmp/commander.yaml

# ... and execute!
$ ./commander test
Starting test file commander.yaml...

✓ echo hello

Duration: 0.002s
Count: 1, Failed: 0
```

## Complete YAML file

Here you can see an example with all features.

```yaml
config: # Config for all executed tests
    dir: /tmp #Set working directory
    env: # Environment variables
        KEY: global
    timeout: 50s # Define a timeout for a command under test
    retries: 2 # Define retries for each test
    
tests:
    echo hello: # Define command as title
        stdout: hello # Default is to check if it contains the given characters
        exit-code: 0 # Assert exit-code
        
    it should fail:
        command: invalid
        stderr:
            contains: 
                - invalid # Assert only contain work
            not-contains:
                - not in there # Validate that a string does not occur in stdout
            exactly: "/bin/sh: 1: invalid: not found"
            line-count: 1 # Assert amount of lines
            lines: # Assert specific lines
                1: "/bin/sh: 1: invalid: not found"
        exit-code: 127
        
    it has configs:
        command: echo hello
        stdout:
            contains: 
                - hello #See test "it should fail"
            exactly: hello
            line-count: 1
        config:
            dir: /home/user # Overwrite working dir
            env:
                KEY: local # Overwrite env variable
                ANOTHER: yeah # Add another env variable
            timeout: 1s # Overwrite timeout
            retries: 5
        exit-code: 0
```

## Executing

```bash
# Execute file commander.yaml in current directory
$ ./commander test 

# Execute a specific suite
$ ./commander test /tmp/test.yaml

# Execute a single test
$ ./commander test /tmp/test.yaml "my test"
```

## Adding tests

You can use the `add` argument if you want to `commander` to create your tests.

```bash
# Add a test to the default commander.yaml
$ ./commander add echo hello
written to /tmp/commander.yaml

# Write to a given file
$ ./commander add --file=test.yaml echo hello
written to test.yaml

# Write to stdout and file
$ ./commander add --stdout echo hello
tests:
  echo hello:
    exit-code: 0
    stdout: hello

written to /tmp/commander.yaml

# Only to stdout
$ ./commander add --stdout --no-file echo hello
tests:
  echo hello:
    exit-code: 0
    stdout: hello
```

## Usage

```
NAME:
   Commander - CLI app testing

USAGE:
   commander [global options] command [command options] [arguments...]

COMMANDS:
     test     Execute the test suite
     add      Automatically add a test to your test suite
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```


## Development

```
# Initialise dev environment
$ make init

# Build the project binary
$ make build

# Unit tests
$ make test

# Coverage
$ make test-coverage

# Integration tests
$ make integration

# Add depdencies to vendor
$ make deps
```

# Misc

Heavily inspired by [goss](https://github.com/aelsabbahy/goss).

Similar projects:
 - [bats](https://github.com/sstephenson/bats)
 - [icmd](https://godoc.org/gotest.tools/icmd)
 - [testcli](https://github.com/rendon/testcli)
