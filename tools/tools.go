//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/Code-Hex/gqldoc/cmd/gqldoc"
	_ "gotest.tools/gotestsum"
)
