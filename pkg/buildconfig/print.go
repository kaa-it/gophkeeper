// Package buildconfig describe function to print build variables
package buildconfig

import (
	"bytes"
	"fmt"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func PrintBuildInfo() {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "Build version: %s\n", buildVersion)
	_, _ = fmt.Fprintf(&buf, "Build date: %s\n", buildDate)
	_, _ = fmt.Fprintf(&buf, "Build commit: %s\n", buildCommit)
	_, _ = fmt.Print(buf.String())
}

