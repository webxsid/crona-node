package updatecheck

import (
	"strconv"
	"strings"
)

type semver struct {
	major      int
	minor      int
	patch      int
	prerelease string
}

func isNewerVersion(current, latest string) bool {
	currentVersion, ok := parseSemver(current)
	if !ok {
		return false
	}
	latestVersion, ok := parseSemver(latest)
	if !ok {
		return false
	}
	if currentVersion.major != latestVersion.major {
		return latestVersion.major > currentVersion.major
	}
	if currentVersion.minor != latestVersion.minor {
		return latestVersion.minor > currentVersion.minor
	}
	if currentVersion.patch != latestVersion.patch {
		return latestVersion.patch > currentVersion.patch
	}
	if currentVersion.prerelease == latestVersion.prerelease {
		return false
	}
	if currentVersion.prerelease == "" && latestVersion.prerelease != "" {
		return false
	}
	if currentVersion.prerelease != "" && latestVersion.prerelease == "" {
		return true
	}
	return latestVersion.prerelease > currentVersion.prerelease
}

func normalizeVersion(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "v")
	return value
}

func parseSemver(value string) (semver, bool) {
	value = normalizeVersion(value)
	if value == "" {
		return semver{}, false
	}
	var prerelease string
	if idx := strings.IndexByte(value, '-'); idx >= 0 {
		prerelease = value[idx+1:]
		value = value[:idx]
	}
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return semver{}, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return semver{}, false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return semver{}, false
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return semver{}, false
	}
	return semver{major: major, minor: minor, patch: patch, prerelease: prerelease}, true
}
