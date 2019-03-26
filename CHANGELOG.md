# v0.4.0
 - Add flags to `test` command
   - `--verbose` will print more detailed output
   - `--no-color` will discard all colors
   - `--concurrent [int value]` sets the maximum concurrently executed tests in `go routines`  
 - Add default test concurrency to `runtime.NumCPU() * 5`

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