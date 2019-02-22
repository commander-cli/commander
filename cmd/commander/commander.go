package main

import (
    "github.com/SimonBaeumer/commander/pkg"
    "github.com/SimonBaeumer/commander/pkg/runtime"
    "os"
)

func main() {
    suite := commander.ParseYAMLFile(os.Args[1])
    runtime.Start(suite)
}
