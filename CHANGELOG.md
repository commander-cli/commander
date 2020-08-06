# v2.3.0

 - Add property `skip`, adds the ability to skip test cases 

# v2.2.0

 - Move from `github.com/SimonBaeumer` to `github.com/commander-cli`
 - Remove `filter` as an argument, instead use `--filter` with regex support

# v2.1.0

 - Add registry authentication for `docker` nodes
 - Add environment support for `Node` configurations
 - Add warning if `pass` is set in `Node` as clear text instead of an env variable
 - Added `--dir` flag to execute all commander test suites within a directory. This feature is not recursive.

# v2.0.0
 
  - Added `nodes` which allow remote execution of tests
  - Added `SSHExecutor` and `LocalExecutor`
  - Removed `concurrent` argument from `test` command
  - Added `DockerExecutor`

# v1.3.0

 - Added `xml` assertion to `stdout` and `stderr`
 - Added `json` assertion to `stdout` and `stderr`
 - Remove `cmd` pkg and use `github.com/SimonBaeumer/cmd@v1.1.0` instead
 - Add `inherit-env` to `config`, it enables inheriting the parent's env variables.

# v1.2.2

 - Rename windows binary `commander-windows-386` to `commander-windows-386.exe`
 - Rename windows binary `commander-windows-amd64` to `commander-windows-amd64.exe`
 - Use commander v1.2.1 in travis build

# v1.2.1

 - Fix `add` command if `stdout` or `stderr` properties were removed if a new test was added

# v1.2.0

 - Add reading environment variables from shell, i.e. `${PATH}`
 - Add `interval` option for `retries` which allows to execute a retry after a given period of time. I.e. `interval: 50ms`

# v1.1.0

 - Add `not-contains` assertion on `stdout` and `stderr`
 - Add validation for invalid data types in `stdout` and `stderr` assertions
 - More logging for `--verbose` option on the `test` command
 - Add better diff format for `contains` and `not-contains` assertions on `stdout` and `stderr`

# v1.0.1

 - Remove unnecessary command logs

# v1.0.0

 - Add `add` command which automatically adds tests to your test suite
 - Changes to `config` properties
    - Add `retries` to test configs
    - Add time units to `timeout` config

# v0.4.0

 - Add flags to `test` command
   - `--verbose` will print more detailed output
   - `--no-color` will discard all colors
   - `--concurrent [int value]` sets the maximum concurrently executed tests in `go routines`  
 - Add default test concurrency to `runtime.NumCPU() * 5`
 - Add more details to log output for each test if `--verbose` is set

# v0.3.0

 - Add `windows` release
 - Add `darwin-386` release
 - Start counting for `lines` in `Stdout` and `Stderr` at `1` instead of `0`
 - Use `maps` instead of `slices` for env variable

# v0.2.1

 - Add `darwin` release

# v0.2.0
 
 - Add test configurations
    - Add the possibility to define environment variables for commands.
    - Add the possibility to set the current working directory for a tested command.
    - Add field validation. If a field does not exist, i.e. due to a typo, display an error message.
    - Add `timeout` to command config. Define a `timeout` in `ms` after which a executed command should fail.
 - Print more error details if a test fails.

# v0.1.0

 - Initial release