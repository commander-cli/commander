[![Build Status](https://travis-ci.org/SimonBaeumer/cmd.svg?branch=master)](https://travis-ci.org/SimonBaeumer/cmd)
[![GoDoc](https://godoc.org/github.com/SimonBaeumer/cmd?status.svg)](https://godoc.org/github.com/SimonBaeumer/cmd)
[![Test Coverage](https://api.codeclimate.com/v1/badges/af3487439a313d580619/test_coverage)](https://codeclimate.com/github/SimonBaeumer/cmd/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/af3487439a313d580619/maintainability)](https://codeclimate.com/github/SimonBaeumer/cmd/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonBaeumer/cmd)](https://goreportcard.com/report/github.com/SimonBaeumer/cmd)

# cmd package

A simple package to execute shell commands on linux, darwin and windows.

## Installation

`$ go get -u github.com/SimonBaeumer/cmd@v1.0.0`

## Usage

```go
c := cmd.NewCommand("echo hello")

err := c.Execute()
if err != nil {
    panic(err.Error())    
}

fmt.Println(c.Stdout())
fmt.Println(c.Stderr())
```

### Configure the command

To configure the command a option function will be passed which receives the command object as an argument passed by reference.

Default option functions:

 - `cmd.WithStandardStreams`
 - `cmd.WithTimeout(time.Duration)`
 - `cmd.WithoutTimeout`
 - `cmd.WithWorkingDir(string)`

#### Example

```go
c := cmd.NewCommand("echo hello", cmd.WithStandardStreams)
c.Execute()
```

#### Set custom options

```go
setWorkingDir := func (c *Command) {
    c.WorkingDir = "/tmp/test"
}

c := cmd.NewCommand("pwd", setWorkingDir)
c.Execute()
```

## Development

### Running tests

```
make test
```

### ToDo

 - os.Stdout and os.Stderr output access after execution via `c.Stdout()` and `c.Stderr()`
