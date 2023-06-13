// Package versioninfo uses runtime.ReadBuildInfo() to set global executable revision information if possible.
package version

import (
	"fmt"
	"time"
)

var (
	// ModulePath is the module path of the application.
	ModulePath = "unknown"
	// Version will be the version tag if the binary is built with "go install url/tool@version".
	// If the binary is built some other way, it will be "(devel)".
	Version = "unknown"
	// Revision is taken from the vcs.revision tag in Go 1.18+.
	Revision = "unknown"
	// LastCommit is taken from the vcs.time tag in Go 1.18+.
	LastCommit time.Time
	// DirtyBuild is taken from the vcs.modified tag in Go 1.18+.
	DirtyBuild = true
)

func String() string {
	return fmt.Sprintf(
		"Module: %v\nVersion: %v\nRevision: %v\nCommitted: %v\nDirty: %v",
		ModulePath,
		Version,
		Revision,
		LastCommit,
		DirtyBuild,
	)
}
