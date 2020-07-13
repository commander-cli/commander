[![Build Status](https://travis-ci.org/commander-cli/cmd.svg?branch=master)](https://travis-ci.org/commander-cli/cmd)
[![GoDoc](https://godoc.org/github.com/commander-cli/cmd?status.svg)](https://godoc.org/github.com/commander-cli/cmd)
[![Test Coverage](https://api.codeclimate.com/v1/badges/af3487439a313d580619/test_coverage)](https://codeclimate.com/github/commander-cli/cmd/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/af3487439a313d580619/maintainability)](https://codeclimate.com/github/commander-cli/cmd/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/commander-cli/cmd)](https://goreportcard.com/report/github.com/commander-cli/cmd)

# cmd package

A simple package to execute shell commands on linux, darwin and windows.

## Installation

`$ go get -u github.com/commander-cli/cmd@v1.0.0`

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
 - `cmd.WithCustomStdout(...io.Writers)`
 - `cmd.WithCustomStderr(...io.Writers)`
 - `cmd.WithTimeout(time.Duration)`
 - `cmd.WithoutTimeout`
 - `cmd.WithWorkingDir(string)`
 - `cmd.WithEnvironmentVariables(cmd.EnvVars)`
 - `cmd.WithInheritedEnvironment(cmd.EnvVars)`

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

### Testing

You can catch output streams to `stdout` and `stderr` with `cmd.CaptureStandardOut`. 

```golang
// caputred is the captured output from all executed source code
// fnResult contains the result of the executed function
captured, fnResult := cmd.CaptureStandardOut(func() interface{} {
    c := NewCommand("echo hello", cmd.WithStandardStream)
    err := c.Execute()
    return err
})

// prints "hello"
fmt.Println(captured)
```

## Development

### Running tests

```
make test
```

### ToDo

 - os.Stdout and os.Stderr output access after execution via `c.Stdout()` and `c.Stderr()`
