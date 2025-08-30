package util

import (
	"fmt"
	"runtime/debug"
)

var (
	Version    = "dev"
	CommitHash = ""
)

func GetVersion() string {
	info, ok := debug.ReadBuildInfo()

	if info.Main.Version != "" && ok {
		return info.Main.Version
	} else if Version == "dev" {
		return fmt.Sprintf("%s-%s", Version, CommitHash)
	} else {
		return Version
	}
}
