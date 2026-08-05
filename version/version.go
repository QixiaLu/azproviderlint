// Package version records the version and git commit the azproviderlint binary was built from.
package version

import "runtime/debug"

var Version = "dev"

var GitCommit string

func init() {
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" && info.Main.Version != "" {
			Version = info.Main.Version
		}
	}
}
