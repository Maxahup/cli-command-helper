package runner

import (
	"os"
	"runtime"
	"strings"
)

type DetectedContext struct {
	OS    string
	Shell string
}

func GetCurrentContext() DetectedContext {
	osName := runtime.GOOS

	shellPath := os.Getenv("SHELL")
	shellName := "unknown"

	if shellPath != "" {
		parts := strings.Split(shellPath, "/")
		shellName = parts[len(parts)-1]
	}
	if osName == "windows" && shellName == "unknown" {
		if os.Getenv("PSModulePath") != "" {
			shellName = "powershell"
		} else {
			shellName = "cmd"
		}
	}

	return DetectedContext{
		OS:    osName,
		Shell: shellName,
	}
}
