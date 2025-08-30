package util

import "fmt"

var (
	Version    = "dev"
	CommitHash = ""
)

func GetVersion() string {
	if Version == "dev" {
		return fmt.Sprintf("%s-%s", Version, CommitHash)
	} else {
		return Version
	}
}
