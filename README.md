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
 - go api
 - logging / verbose output
 - command execution
   - environment variables
   - arguments?
   - timeout
 - exit code
 - stdout
    - Validate against string
    - Validate against file
    - Validate against line
    - Validate with wildcards / regex
 - stderr
    - Validate against string
    - Validate against file
    - Validate with wildcards
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
