package application

import (
	"time"
)

// For backward compatibility.
const fallbackTimeLayout = "2006-01-02 15:04:05"

type BuildVersion struct {
	Time   time.Time
	Commit string
}

// These values have to be set via LDFLAGS.
var (
	buildVersionTime   string
	buildVersionCommit string
)

// buildVersion - stores parsed value.
var buildVersion BuildVersion

func GetBuildVersion() (BuildVersion, error) {
	if !buildVersion.Time.IsZero() {
		return buildVersion, nil
	}
	buildVersion.Commit = buildVersionCommit

	var err error

	buildVersion.Time = time.Now()

	return buildVersion, err
}
