package main

import (
	"fmt"
	"os"

	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/katbyte/azproviderlint/checks"
	"github.com/katbyte/azproviderlint/version"
)

func main() {
	// multichecker owns the CLI, so handle `azproviderlint version` before it takes over
	if len(os.Args) > 1 && os.Args[1] == "version" {
		v := version.Version
		if v == "" {
			v = "dev"
		}
		if version.GitCommit != "" {
			v += " (" + version.GitCommit + ")"
		}
		fmt.Println("azproviderlint " + v)
		return
	}

	multichecker.Main(checks.All...)
}
