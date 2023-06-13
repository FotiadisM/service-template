//go:build go1.12

package version

import "runtime/debug"

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	Version = info.Main.Version
	ModulePath = info.Path
}
