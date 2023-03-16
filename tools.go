//go:build tools
// +build tools

// see https://play-with-go.dev/tools-as-dependencies_go115_en/
// and https://www.jvt.me/posts/2022/06/15/go-tools-dependency-management/
package tools

import (
	_ "mvdan.cc/gofumpt"
)
