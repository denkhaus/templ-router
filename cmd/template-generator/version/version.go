package version

import (
	"fmt"
	"runtime/debug"
)

// These will be set at build time via ldflags
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// BuildInfo contains version and build information
type BuildInfo struct {
	Version   string
	GitCommit string
	BuildTime string
	GoVersion string
	Platform  string
}

// GetBuildInfo returns comprehensive build information
func GetBuildInfo() BuildInfo {
	info := BuildInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		GoVersion: fmt.Sprintf("%s", getGoVersion()),
		Platform:  fmt.Sprintf("%s/%s", getGOOS(), getGOARCH()),
	}

	// Try to get more info from build info if available
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		// Get Go version from build info if not set
		if info.GoVersion == "" {
			info.GoVersion = buildInfo.GoVersion
		}

		// Try to get commit info from VCS if not set via ldflags
		if info.GitCommit == "unknown" {
			for _, setting := range buildInfo.Settings {
				switch setting.Key {
				case "vcs.revision":
					if len(setting.Value) >= 7 {
						info.GitCommit = setting.Value[:7] // Short commit hash
					} else {
						info.GitCommit = setting.Value
					}
				case "vcs.time":
					if info.BuildTime == "unknown" {
						info.BuildTime = setting.Value
					}
				}
			}
		}
	}

	return info
}

// String returns a formatted version string
func (b BuildInfo) String() string {
	return fmt.Sprintf("template-generator %s (commit: %s, built: %s, go: %s, platform: %s)",
		b.Version, b.GitCommit, b.BuildTime, b.GoVersion, b.Platform)
}

// Short returns a short version string
func (b BuildInfo) Short() string {
	return fmt.Sprintf("v%s-%s", b.Version, b.GitCommit)
}

// GetVersion returns the version string
func GetVersion() string {
	return GetBuildInfo().String()
}

// GetShortVersion returns a short version string
func GetShortVersion() string {
	return GetBuildInfo().Short()
}

// Helper functions to get runtime info
func getGoVersion() string {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		return buildInfo.GoVersion
	}
	return "unknown"
}

func getGOOS() string {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range buildInfo.Settings {
			if setting.Key == "GOOS" {
				return setting.Value
			}
		}
	}
	return "unknown"
}

func getGOARCH() string {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range buildInfo.Settings {
			if setting.Key == "GOARCH" {
				return setting.Value
			}
		}
	}
	return "unknown"
}
