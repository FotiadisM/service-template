//go:build tools

package tools

import (
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"

	_ "gotest.tools/gotestsum"
	_ "mvdan.cc/gofumpt"
)
