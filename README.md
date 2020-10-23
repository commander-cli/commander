[![Build Status](https://travis-ci.org/commander-cli/commander.svg?branch=master)](https://travis-ci.org/commander-cli/commander)
[![GoDoc](https://godoc.org/github.com/commander-cli/commander?status.svg)](https://godoc.org/github.com/commander-cli/commander)
[![Go Report Card](https://goreportcard.com/badge/github.com/commander-cli/commander)](https://goreportcard.com/report/github.com/commander-cli/commander)
[![Maintainability](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/maintainability)](https://codeclimate.com/github/commander-cli/commander/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/cc848165784e0f809a51/test_coverage)](https://codeclimate.com/github/commander-cli/commander/test_coverage)
[![Github All Releases](https://img.shields.io/github/downloads/commander-cli/commander/total.svg)](https://github.com/commander-cli/commander/releases)

# Commander

Define language independent tests for your command line scripts and programs in simple `yaml` files.

 - It runs on `windows`, `osx` and `linux` 
 - It can validate local machines, ssh hosts and docker containers
 - It is a self-contained binary - no need to install a heavy lib or language
 - It is easy and fast to write

For more information take a look at the [quick start](#quick-start), the [examples](examples) or the [integration tests](integration).

## Table of contents

* [Installation](#installation)
  + [Any system with Go installed](#any-system-with-go-installed)
  + [Linux & osx](#linux---osx)
  + [Windows](#windows)
* [Quick start](#quick-start)
  + [Complete YAML file](#complete-yaml-file)
  + [Executing](#executing)
  + [Adding tests](#adding-tests)
* [Documentation](#documentation)
  + [Usage](#usage)
  + [Tests](#tests)
    - [command](#command)
    - [config](#user-content-config-test)
    - [exit-code](#exit-code)
    - [stdout](#stdout)
      * [contains](#contains)
      * [exactly](#exactly)
      * [json](#json)
      * [lines](#lines)
      * [line-count](#line-count)
      * [not-contains](#not-contains)
      * [xml](#xml)
      * [file](#file)
    - [stderr](#stderr)
    - [skip](#skip)
  + [Config](#user-content-config-config)
    - [dir](#dir)
    - [env](#env)
    - [inherit-env](#inherit-env)
    - [interval](#interval)
    - [retries](#retries)
    - [timeout](#timeout)
    - [nodes](#nodes)
  + [Nodes](#nodes)
    - [local](#local)
    - [ssh](#ssh)
    - [docker](#docker)
  + [Development](#development)
* [Misc](#misc)

## Installation

### Any system with Go installed

Probably the easiest way to install `commander` is by using `go get` to download and install it in one simple command:

```bash
go get github.com/commander-cli/commander/cmd/commander
```

This works on any OS, as long as go is installed. If go is not installed on your system, follow one of the methods below.

### Linux & osx

Visit the [release](https://github.com/commander-cli/commander/releases) page to get the binary for you system. 

```bash
curl -L https://github.com/commander-cli/commander/releases/download/v1.2.2/commander-linux-amd64 -o commander
chmod +x commander
```

### Windows

 - Download the current [release](https://github.com/commander-cli/commander/releases/latest)
 - Add the path to your [path](https://docs.alfresco.com/4.2/tasks/fot-addpath.html) environment variable
 - Test it: `commander --version`


## Quick start

A `commander` test suite consists of a `config` and `tests` root element. To start quickly you can use 
the following examples.

```bash
# You can even let commander add tests for you!
$ ./commander add --stdout --file=/tmp/commander.yaml echo hello
tests:
  echo hello:
    exit-code: 0
    stdout: hello

written to /tmp/commander.yaml

# ... and execute!
$ ./commander test /tmp/commander.yaml
Starting test file /tmp/commander.yaml...

✓ echo hello

Duration: 0.002s
Count: 1, Failed: 0
```

### Complete YAML file

Here you can see an example with all features for a quick reference

```yaml
nodes:
  ssh-host1:
    type: ssh
    addr: 192.168.0.1:22
    user: root
    pass: pass
  ssh-host2:
    type: ssh
    addr: 192.168.0.1:22
    user: root
    identity-file: /home/user/id_rsa.pub
  docker-host1:
    type: docker
    image: alpine:2.4
  docker-host2:
    type: docker
    instance: alpine_instance_1

config: # Config for all executed tests
    dir: /tmp #Set working directory
    env: # Environment variables
        KEY: global
    timeout: 50s # Define a timeout for a command under test
    retries: 2 # Define retries for each test
    nodes:
    - ssh-host1 # define default hosts
    
tests:
    echo hello: # Define command as title
        stdout: hello # Default is to check if it contains the given characters
        exit-code: 0 # Assert exit-code
        
    it should skip:
        command: echo "I should be skipped"
        stdout: I should be skipped
        skip: true
        
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
            json:
                object.attr: hello # Make assertions on json objects
            xml:
                "//book//auhtor": Steven King # Make assertions on xml documents
            file: correct-output.txt
        exit-code: 127
        skip: false

    it has configs:
        command: echo hello
        stdout:
            contains: 
              - hello #See test "it should fail"
            exactly: hello
            line-count: 1
        config:
            inherit-env: true # You can inherit the parent shells env variables
            dir: /home/user # Overwrite working dir
            env:
                KEY: local # Overwrite env variable
                ANOTHER: yeah # Add another env variable
            timeout: 1s # Overwrite timeout
            retries: 5
            nodes: # overwrite default nodes
              - docker-host1
              - docker-host2
        exit-code: 0
```

### Executing

```bash
# Execute file commander.yaml in current directory
$ ./commander test 

# Execute a specific suite
$ ./commander test /tmp/test.yaml

# Execute a single test
$ ./commander test /tmp/test.yaml --filter "my test"

# Execute suite from stdin
$ cat /tmp/test.yaml | ./commander test -

# Execute suite from url
$ ./commander test https://your-url/commander_test.yaml

# Execute suites within a test directory
$ ./commander test --dir /tmp
```

### Adding tests

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

## Documentation

### Usage

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


### Tests

Tests are defined in the `tests` root element. Every test consists of a [command](#command) and an expected result, 
i.e. an [exit-code](#exit-code). 


```yaml
tests: # root element
  echo test: # test case - can either be the command or a given title
    stdout: test
    exit-code: 0
```

A test is a `map` which configures the test. 
The `key` (`echo test` in the example above) of the test can either be the `command` itself or the `title` of the test which will be displayed in the test execution.

If the same `command` is tested multiple times it is useful to set the `title` of the test manually and use the `command` property. 
Further the `title` can be useful to describe tests better. See [the commander test suite](commander_unix.yaml) as an example.

 - name: `title or command under test`
 - type: `map`
 - default: `{}`
 
Examples:

```yaml
tests:
  echo test: # command and title will be the same
    stdout: test
    exit-code: 0
    
  my title: # custom title
    command: exit 1 # set command manually
    exit-code: 1
```

#### command

`command` is a `string` containing the `command` to be tested. Further the `command` property is automatically parsed from 
the `key` if no `command` property was given.

 - name: `command`
 - type: `string`
 - default: `can't be empty`
 - notes: Will be parsed from the `key` if no `command` property was provided and used as the title too
 
 
```yaml
echo test: # use command as key and title
  exit-code: 0
  
it should print hello world: # use a more descriptive title...
  command: echo hello world  # ... and set the command in the property manually
  stdout: hello world
  exit-code: 0
```

#### <a name="config-test"></a>config

`config` sets configuration for the test. `config` can overwrite global configurations. 

 - name: `config`
 - type: `map`
 - default: `{}`
 - notes:
   - for more information look at [config](#user-content-config-config)

```yaml
echo test:
  config:
    timeout: 5s
```

#### exit-code

`exit-code` is an `int` type and compares the given code to the `exit-code` of the given command.

 - name: `exit-code`
 - type: `int`
 - default: `0`
 
```yaml
exit 1: # will pass
  exit-code: 1
exit 0: # will fail
  exit-code: 1
```

#### stdout

`stdout` and `stderr` allow to make assertions on the output of the command. 
The type can either be a `string` or a `map` of different assertions.

If only a `string` is provided it will check if the given string is [contained](#contains) in the output.

 - name: `stdout`
 - type: `string` or `map`
 - default: ` `
 - notes: [stderr](#stderr) works the same way
 
```yaml
echo test:
  stdout: test # make a contains assertion
  
echo hello world:
  stdout:
    line-count: 1 # assert the amount of lines and use stdout as a map
```

##### contains

`contains` is an `array` or `string`. It checks if a `string` is contained in the output. 
It is the default if a `string` is directly assigned to `stdout` or `stderr`.

 - name: `contains`
 - type: `string` or `array`
 - default: `[]`
 - notes: default assertion if directly assigned to `stdout` or `stderr`

```yaml
echo hello world:
  stdout: hello # Default is a contains assertion

echo more output:
  stdout:
    contains:
      - more
      - output
```

##### exactly

`exactly` is a `string` type which matches the exact output.

 - name: `exactly`
 - type: `string`
 - default: ` `
 
```yaml
echo test:
  stdout:
    exactly: test
```

##### json

`json` is a `map` type and allows to parse `json` documents with a given `GJSON syntax` to query for specific data. 
The `key` represents the query, the `value` the expected value.

 - name: `json`
 - type: `map`
 - default: `{}`
 - notes: Syntax taken from [GJSON](https://github.com/tidwall/gjson#path-syntax)
 
```yaml
cat some.json: # print json file to stdout
  name.last: Anderson # assert on name.last, see document below
``` 

`some.json` file:
 
```json
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}
```

More examples queries:

```
"name.last"          >> "Anderson"
"age"                >> 37
"children"           >> ["Sara","Alex","Jack"]
"children.#"         >> 3
"children.1"         >> "Alex"
"child*.2"           >> "Jack"
"c?ildren.0"         >> "Sara"
"fav\.movie"         >> "Deer Hunter"
"friends.#.first"    >> ["Dale","Roger","Jane"]
"friends.1.last"     >> "Craig"
```

##### lines

`lines` is a `map` which makes exact assertions on a given line by line number.

 - name: `lines`
 - type: `map`
 - default: `{}`
 - note: starts counting at `1` ;-)
 
```yaml
echo test\nline 2:
  stdout:
    lines:
      2: line 2 # asserts only the second line
```

##### line-count

`line-count` asserts the amount of lines printed to the output. If set to `0` this property is ignored.

 - name: `line-count`
 - type: `int`
 - default: `0`

```yaml
echo test\nline 2:
  stdout:
    line-count: 2
```

##### not-contains

`not-contains` is a `array` of elements which are not allowed to be contained in the output. 
It is the inversion of [contains](#contains).

 - name: `not-contains`
 - type: `array`
 - default: `[]`

```yaml
echo hello:
  stdout:
    not-contains: bonjour # test passes because bonjour does not occur in the output
 
echo bonjour:
  stdout:
    not-contains: bonjour # test fails because bonjour occurs in the output
```

##### xml

`xml` is a `map` which allows to query `xml` documents viá `xpath` queries. 
Like the [json][#json] assertion this uses the `key` of the map as the query parameter to, the `value` is the expected value.

 - name: `xml`
 - type: `map`
 - default: `{}`
 - notes: Used library [xmlquery](https://github.com/antchfx/xmlquery)

```yaml
cat some.xml:
  stdout:
    xml:
      //book//author: J. R. R. Tolkien
```

`some.xml` file:

```xml
<book>
    <author>J. R. R. Tolkien</author>
</book>
```

##### file

`file` is a file path, relative to the working directory that will have
its entire contents matched against the command output. Other than
reading from a file this works the same as [exactly](#exactly).

The example below will always pass.

```yaml
output should match file:
  command: cat output.txt
  stdout:
    file: output.txt
```

#### stderr

See [stdout](#stdout) for more information.

 - name: `stderr`
 - type: `string` or `map`
 - default: ` `
 - notes: is identical to [stdout](#stdout) 

```yaml
# >&2 echos directly to stderr
">&2 echo error": 
  stderr: error
  exit-code: 0

">&2 echo more errors":
  stderr:
    line-count: 1
```

#### skip

`skip` is a `boolean` type, setting this field to `true` will skip the test case.

 - name: `skip`
 - type: `bool`
 - default: `false`

```yaml
echo test:
  stdout: test
  skip: true
```

### <a name="config-config"></a>Config

You can add configs which will be applied globally to all tests or just for a specific test case, i.e.:

```yaml
config:
  dir: /home/root # Set working directory for all tests

tests:
  echo hello:
    config: # Define test specific configs which overwrite global configs
      timeout: 5s
  exit-code: 0
```

#### dir

`dir` is a `string` which sets the current working directory for the command under test. 
The test will fail if the given directory does not exist.

 - name: `dir`
 - type: `string`
 - default: `current working dir`

```yaml
dir: /home/root
```

#### env

`env` is a `hash-map` which is used to set custom env variables. The `key` represents the variable name and the `value` setting the value of the env variable.

 - name: `env`  
 - type: `hash-map`
 - default: `{}`
 - notes:
    - read env variables with `${PATH}`
    - overwrites inherited variables, see [#inherit-env](#inherit-env) 
    
```yaml
env:
  VAR_NAME: my value # Set custom env var
  CURRENT_USER: ${USER} # Set env var and read from current env
```

#### inherit-env

`inherit-env` is a `boolean` type which allows you to inherit all environment variables from your active shell.

 - name: `inherit-env`
 - type: `bool`
 - default: `false`
 - notes: If this config is set to `true` in the global configuration it will be applied for all tests and ignores local test configs.

```yaml
inherit-env: true
```

#### interval

`interval` is a `string` type and sets the `interval` between [retries](#retries).
 
 - name: `interval`
 - type: `string`
 - default: `0ns`
 - notes:
   - valid time units: ns, us, µs, ms, s, m, h
   - time string will be evaluated by golang's `time` package, further reading [time/#ParseDuration](https://golang.org/pkg/time/#ParseDuration)

```yaml
interval: 5s # Waits 5 seconds until the next try after a failed test is started
```

#### retries

`retries` is an `int` type and configures how often a test is allowed to fail until it will be marked as failed for the whole test run.

 - name: `retries`
 - type: `int`
 - default: `0`
 - notes: [interval](#interval) can be defined between retry executions

```yaml
retries: 3 # Test will be executed 3 times or until it succeeds
```

#### timeout

`timeout` is a `string` type and sets the time a test is allowed to run. 
The time is parsed from a duration string like `300ms`.
If a tests exceeds the given `timeout` the test will fail.

 - name: `timeout`
 - type: `string`
 - default: `no limit`
 - notes:
   - valid time units: ns, us, µs, ms, s, m, h
   - time string will be evaluated by golang's `time` package, further reading [time/#ParseDuration](https://golang.org/pkg/time/#ParseDuration)

```yaml
timeout: 600s
```

### Nodes

`Commander` has the option to execute tests against other hosts, i.e. via ssh.

Available node types are currently:

 - `local`, execute tests locally
 - `ssh`, execute tests viá ssh
 - `docker`, execute tests inside a docker container

```yaml
nodes: # define nodes in the node section
  ssh-host:
    type: ssh # define the type of the connection 
    user: root # set the user which is used by the connection
    pass: password # set password for authentication
    addr: 192.168.0.100:2222 # target host address
    identity-file: ~/.ssh/id_rsa # auth with private key
tests:
  echo hello:
    config:
      nodes: # define on which host the test should be executed
        - ssh-host
    stdout: hello
    exit-code: 0
```

You can identify on which node a test failed by inspecting the test output.
The `[local]` and `[ssh-host]` represent the node name on which the test were executed.

```
✗ [local] it should test ssh host
✗ [ssh-host] it should fail if env could not be set
```

#### local

The `local` node is the default execution and will be applied if nothing else was configured.
It is always pre-configured and available, i.e. if you want to execute tests on a node and locally.

```yaml
nodes:
  ssh-host:
    addr: 192.168.1.100
    user: ...
tests:
  echo hello:
    config:
      nodes: # will be executed on local and ssh-host
        - ssh-host
        - local
    exit-code: 0
```

#### ssh

The `ssh` node type will execute tests against a configured node using ssh.

**Limitations:** 
 - The `inhereit-env` config is disabled for ssh hosts, nevertheless it is possible to set env variables
 - Private registries are not supported at the moment

```yaml
nodes: # define nodes in the node section
  ssh-host:
    type: ssh # define the type of the connection 
    user: root # set the user which is used by the connection
    pass: password # set password for authentication
    addr: 192.168.0.100:2222 # target host address
    identity-file: ~/.ssh/id_rsa # auth with private key
tests:
  echo hello:
    config:
      nodes: # define on which host the test should be executed
        - ssh-host
    stdout: hello
    exit-code: 0
```

#### docker

The `docker` node type executes the given command inside a docker container.

**Notes:** If the default docker registry should be used prefix the container with the registry `docker.io/library/` 

```yaml:
nodes:
  docker-host:
    type: docker
    image: docker.io/library/alpine:3.11.3
    docker-exec-user: 1000 # define the owner of the executed command
    user: user # registry user
    pass: password # registry password, it is recommended to use env variables like $REGISTRY_PASS
config:
  nodes:
    - docker-host
    
tests:
  "id -u":
     stdout: "1001"
```

### Development

See the documentation at [development.md](docs/development.md)

## Misc

Heavily inspired by [goss](https://github.com/aelsabbahy/goss).

Similar projects:
 - [bats](https://github.com/sstephenson/bats)
 - [icmd](https://godoc.org/gotest.tools/icmd)
 - [testcli](https://github.com/rendon/testcli)
