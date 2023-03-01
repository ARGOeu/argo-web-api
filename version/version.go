package version

import (
	"runtime"
)

var (
	// Release version of the service. Bump it up during new version release
	Release = "1.13.2"
	// Commit hash provided during build
	Commit = "Unknown"
	// BuildTime provided during build
	BuildTime = "Unknown"
	// GO provides golang version
	GO = runtime.Version()
	// Compiler info
	Compiler = runtime.Compiler
	// OS Info
	OS = runtime.GOOS
	// Arch info
	Arch = runtime.GOARCH
)
