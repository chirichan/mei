package version

import (
	"runtime/debug"
	"strings"
)

var (
	Version   string = "v0.1.7"
	BuildTime string
	GitCommit string
	GoVersion string
)

const unknown = "unknown"

func ResolvedVersion() string {
	buildVersion := buildInfoMainVersion()
	if buildVersion != "" {
		return buildVersion
	}
	if Version != "" {
		return Version
	}
	return unknown
}

func ResolvedGitCommit() string {
	if GitCommit != "" {
		return GitCommit
	}
	revision := buildInfoSetting("vcs.revision")
	if revision == "" {
		return unknown
	}
	if len(revision) > 8 {
		return revision[:8]
	}
	return revision
}

func ResolvedBuildTime() string {
	if BuildTime != "" {
		return BuildTime
	}
	vcsTime := buildInfoSetting("vcs.time")
	if vcsTime != "" {
		return vcsTime
	}
	return unknown
}

func buildInfoMainVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	mainVersion := strings.TrimSpace(info.Main.Version)
	if mainVersion == "" || mainVersion == "(devel)" {
		return ""
	}
	return mainVersion
}

func buildInfoSetting(key string) string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, s := range info.Settings {
		if s.Key == key {
			return strings.TrimSpace(s.Value)
		}
	}
	return ""
}
