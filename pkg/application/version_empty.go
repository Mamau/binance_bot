//go:build empty_build_version

package application

import "time"

func init() {
	buildVersionTime = time.Time{}.Format(fallbackTimeLayout)
	buildVersionCommit = "none"
}
