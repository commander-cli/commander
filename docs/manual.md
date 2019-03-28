# Manual

`commander` will automatically search for a `commander.yaml` in the current working directory.

## Commands

## add

The `add` command allows you to automatically create tests. It will automatically create a `comannder.yaml`.

**Example**

```bash
./commander add --no-file --stdout echo hello
tests:
  echo hello:
    exit-code: 0
    stdout: hello
```

**Options**

```
--stdout      Output test file to stdout
--no-file     Don't create a commander.yaml
--file value  Write to another file, default is commander.yaml
```

## test

```yaml
config: # Config for all tests
    dir: /tmp #Set working directory
    env: # Environment variables
        KEY: global
    timeout: 5000 #Timeout in ms
    
tests:
    echo hello: # Define command as title
        stdout: hello # Default is to check if it contains the given characters
        exit-code: 0 # Assert exit-code
        
    it should fail:
        command: invalid
        stderr:
            contains: 
                - invalid # Assert only contain work
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
            timeout: 1000 # Overwrite timeout
        exit-code: 0
```