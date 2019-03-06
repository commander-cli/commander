[![Build Status](https://travis-ci.org/SimonBaeumer/commander.svg?branch=master)](https://travis-ci.org/SimonBaeumer/commander)
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonBaeumer/commander)](https://goreportcard.com/report/github.com/SimonBaeumer/commander)
[![Maintainability](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/maintainability)](https://codeclimate.com/github/SimonBaeumer/commander/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/test_coverage)](https://codeclimate.com/github/SimonBaeumer/commander/test_coverage)
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
 - assert stderr & exit-code & line-count
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

## Matching / Printing

- Diff libs
    - https://github.com/sergi/go-diff
    - https://github.com/pmezard/go-difflib
    - https://github.com/sourcegraph/go-diff
- regex matching
- line count
- assert specific lines
- specific line

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
        skip: true
        exit-code: 0
```

## Docker

```
tests:
    echo hello:
        docker: ubuntu:18.04
        exit-code: 0
```