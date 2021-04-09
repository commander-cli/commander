# Development documentation

- [Introduction](#introduction)
  * [Directory overview](#directory-overview)
  * [Package overview](#package-overview)
  * [Build targets](#build-targets)
  * [Unit tests](#unit-tests)
- [Extending commander](#extending-commander)
  * [Add a new field to the `YAML` suite - with a leaning-by-doing task](#add-a-new-field-to-the--yaml--suite)
  * [Writing integration tests](#writing-integration-tests)


## Introduction

### Directory overview

```
├── cmd          # Contains the composition root and files which will be compiled to binaries.
├── docs         # Documentation for commander
├── examples     # Examples of how to use commander with to give you a little bit inspiration.
├── hack         # Just a directory for testing stuff and shitty dev scripts.
├── integration  # Integration tests for all platforms, all written as test suites for commander.
├── pkg          # All packages written in GoLang which are used to compose this tool.
├── release      # All binaries which are created on `make release` are located here. Directory is in `.gitignore`.
└── vendor       # Third party dependencies added by `go mod`.
```

### Package overview

| **package** | **description**                                                                                                                                                                                                                                                                                                 | **path**       |
|-------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------|
| main        | The composition root of the commander, initializes the `urfave/cli` framework.                                                                                                                                                                                                                                  | /cmd/commander |
| app         | Contains all commands for commander, i.e. add or test. Will only be used by the `cmd` package located under `/cmd`                                                                                                                                                                                              | /pkg/app       |
| matcher     | Implements the matching logic for the different assertions types in the suite like `equals, `contains`, `json` or `xml`.                                                                                                                                                                                        | /pkg/matcher   |
| output      | Output is a package which allows to add different output types. For example it could be possible to add a new output format in `json` or for a health check.                                                                                                                                                    | /pkg/output    |
| runtime     | Runtime controls the execution of the test suite. It will be initialized by the `main` package and will be a given a `Suite` which contains all tests to be executed. The runtime also contains all `executors` which are different types of ways to execute a test, i.e. on a local machine or a node viá ssh. | /pkg/runtime   |
| suite       | Suite is responsible for parsing the defined formats of different test suites. At the moment only `yaml` is supported but it would be possible to add support for a custom DSL or formats like `toml` and `json`. The `suite.Suite` struct will be used by the `runtime.Runtime` for executing the tests.       | /pkg/suite     |


### Build targets

```
# Initialise dev environment
$ make init

# Build the project binary
$ make build

# Unit tests
$ make test

# Coverage
$ make test-coverage

# Coverage with more complex tests like ssh execution
$ make test-coverage-all-dockerized

# Integration tests for linux and macos
$ make integration-unix

# Integration on linux
$ make integration-linux

# Integration windows
$ make integration-windows

# Add depdencies to vendor
$ make deps
```

### Unit tests

`COMMANDER_TEST_ALL` will enable all tests which are depending on external systems like docker or databases.
Enables ssh tests in unit test suite and sets the credentials for the target host.
`COMMANDER_SSH_TEST` must be set to `1` to enable ssh tests.

**Note:** I am aware that unit tests should not test external systems nor libraries, but in favour of simplicity and laziness
I created simple tests inside the directory tree.

```bash
export COMMANDER_TEST_ALL=1
export COMMANDER_TEST_SSH=1
export COMMANDER_TEST_SSH_HOST=localhost:2222
export COMMANDER_TEST_SSH_PASS=pass
export COMMANDER_TEST_SSH_USER=root
export COMMANDER_TEST_SSH_IDENTITY_FILE=integration/containers/ssh/.ssh/id_rsa
```

## Extending commander

### Add a new field to the `YAML` suite - with a leaning-by-doing task

It is a little bit annoying to add fields because it will be converted multiple times to keep the
suite format abstracted from the runtime package. 
The idea behind this it to add support for other formats like `json`, `toml` or maybe a custom DSL.

**Definition of done**:

 - Add a property `message` which always display a message while executing a test
   ```yaml
   tests:
     echo hello:
       exit-code: 0
       config:
         message: this is a very special test
   ``` 
 - Support global configurations
 - Create a simple test case

 **1. Extend the conversion struct `suite.YAMLTestConfigConf` in [yaml_suite.go](../pkg/suite/yaml_suite.go) with the `Message` property.**
 
   The structs types with a `Conf` suffix represent the configuration type, the `YAML` prefix the format of the suite.
   The naming makes it a little bit clearer in the code which type is used, i.e. a `runtime.CommandUnderTest` or a `suite.YAMLTestConfigConf`.
    
    ```go
    Message    string            `yaml:"message,omitempty"`
    ```
    
**2. Add the `Message` property as a `string` to the `runtime.CommandUnderTest` struct**
 
   The `runtime.CommandUnderTest` struct in [runtime.go](../pkg/runtime/runtime.go) will be used by the runtime to 
   create the command with all its configs like `env` variables and is used for the test execution.
    
**3. Add a simple test case**
 
  Open [yaml_suite_test.go](../pkg/suite/yaml_suite_test.go) and look for an existing test to add this property or create a new one. 
  I recommend to create a simple test case before adding the properties because it is easier to debug and test. 
  In our example we could extend the `TestYAMLSuite_ShouldPreferLocalTestConfigs` test but will add a new `TestYAMLSuite_Message` test for the simplicity.
  
  ```go
  func TestYamlSuite_Message(t *testing.T) {
  yaml := []byte(`
    tests:
      echo hello:
        exit-code: 0
        config:
          message: "This is a very special test"
  `)
    
   got := ParseYAML(yaml)
    
    assert.Equal(t, "This is a very special test", got.TestCases[0].Command.Message)
  }
  ```
    
  Run the unit tests with `make test`. It should print a result like this:
  
  ```
  --- FAIL: TestYamlSuite_Message (0.00s)
      yaml_suite_test.go:192: 
                  Error Trace:    yaml_suite_test.go:192
                  Error:          Not equal: 
                                  expected: "This is a very special test"
                                  actual  : ""
                                  
                                  Diff:
                                  --- Expected
                                  +++ Actual
                                  @@ -1 +1 @@
                                  -This is a very special test
                                  +
                  Test:           TestYamlSuite_Message
  FAIL
  ```
    
**4. Parse YAML and convert it to `suite.Suite`**
 
  This is a little bit complicated and error prone because it is splitted into three steps:
  
  - Parse YAML file
  - Convert YAML config structs to runtime test structs
  - Assign global configuration
  
  For this take a look at the `pkg/suite/yaml_suite.go:ParseYAML` function which is responsible for parsing the suite.
  
  1. Parse yaml - `err := yaml.UnmarshalStrict(content, &yamlConfig)`
      This line parses the yaml file and will return a `suite.YAMLSuiteConf`, later it will be converted to our structs
      from the `runtime` package.

      `Suite::ParseYAML` will then convert the Unmarshalled `suite.YAMLSuiteConf` 
      into a `suite.Suite`. Navigate to `Suite::NewSuite` and follow the code to
      `s.mergeConfigs`.
      
      Jump into the `Suite::mergeConfigs` in [pkg/suite/suite.go]
      (../pkg/suite/suite.go) and add the `message` property like this:
      
      ```go
      if s.Config.Message == "" {
          s.Config.Message = config.Message
      }
      ```

      as well adding similar logic to `Suite::mergeTestConfigs` 
      [pkg/suite/suite.go](../pkg/suite/suite.go)

      ```go
      if s.TestCases[i].Command.Message == "" {
          s.TestCases[i].Command.Message = s.Config.Message
      }
      ```

  1. Convert test cases - `tests := convertYAMLSuiteConfToTestCases(yamlConfig)`
  
     Jump into the `convertYAMLSuiteConfToTestCases` function and assign the content of our new field to the `runtime.CommandUnderTest` conversion.
     
     ```
      Command: runtime.CommandUnderTest{
          Cmd:        t.Command,
          InheritEnv: t.Config.InheritEnv,
          [...]
          Message:    t.Config.Message,
      },
      ```
      
**5. Add the global config assignment and add the `Message` property**

  Add a new `Message` property to the `runtime.GlobalTestConfig`.
  
  ```go
  type GlobalTestConfig struct {
    Env        map[string]string
    [...]
    Message    string
  }
  ```
  
  And last but not least implement the assignment of the global config. 
  This can be done easily inside the `return` of the `ParseYAML` function statement.

   ```go       
    return Suite{
      TestCases: tests,
      Config: runtime.TestConfig{
        InheritEnv: yamlConfig.Config.InheritEnv,
        Env:        yamlConfig.Config.Env,
        Dir:        yamlConfig.Config.Dir,
        Timeout:    yamlConfig.Config.Timeout,
        Retries:    yamlConfig.Config.Retries,
        Interval:   yamlConfig.Config.Interval,
        Nodes:      yamlConfig.Config.Nodes,
      },
      Nodes: convertNodes(yamlConfig.Nodes),
    }
   ```
  
**6. Run the tests**

  ```bash
  $ make test
  INFO: Starting build test
  go test ./...
  ok      github.com/commander-cli/commander/cmd/commander 0.011s
  ok      github.com/commander-cli/commander/pkg/app       0.014s
  ok      github.com/commander-cli/commander/pkg/matcher   (cached)
  ok      github.com/commander-cli/commander/pkg/output    (cached)
  ok      github.com/commander-cli/commander/pkg/runtime   0.229s
  ok      github.com/commander-cli/commander/pkg/suite     0.008s
  ```
     
**7. Learning by doing**

  Two tasks are missing which you can complete on your own.
  
   - Add the printing of the message (tip: Take a look into the `runtime.go` file)
   - Extend the test case that it is tested that local configs are preferred over global configs (tip: Take a look at the other tests). 
  
### Writing integration tests

Commander tests itself. You can find the integration tests in `commander_unix.yaml` and `commander_windows.yaml`.
More complex scenarios are stored in `integration/`.

It is always necessary to execute the test suite with a stable version of commander and not the current build.

**Tipps:**

 - The working directory is by default the project root, even for tests located inside `integration/*`
 - Execute `commander` inside the `commander_*.yaml` files with a given suite and assert the result which is returned
