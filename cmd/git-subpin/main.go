// Command git-subpin resolves Git submodule pointer conflicts during automated
// merges between environment branches, driven by a declarative policy.
//
// This is a distribution stub: it prints the runtime platform and the action
// inputs it received so the GitHub Action packaging (shim, release branch and
// runner matrix) can be verified end to end before the resolution logic lands.
package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

func main() {
	if err := run(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "git-subpin:", err)
		os.Exit(1)
	}
}

func run(out io.Writer) error {
	fmt.Fprintln(out, "git-subpin distribution stub")
	fmt.Fprintf(out, "platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	for _, name := range []string{"FROM", "TO", "CONFIG"} {
		fmt.Fprintf(out, "input %s: %q\n", name, os.Getenv("INPUT_"+name))
	}
	return nil
}
