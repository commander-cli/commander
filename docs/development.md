# Development documentation

## Writing integration tests

Commander tests itself. You can find the integration tests in `commander_unix.yaml` and `commander_windows.yaml`.
More complex scenarios are stored in `integration/`.

It is always necessary to execute the test suite with a stable version of commander and not the current build.

**Tipps:**

 - The working directory is by default the project root, even for tests located inside `integration/*`
 - Execute `commander` inside the `commander_*.yaml` files with a given suite and assert the result which is returned