package version

import (
	"fmt"
	"github.com/blang/semver"
	"runtime"
	"strings"
)

var version, gitCommit, buildDate string
var platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

type Info struct {
	Version    string
	GitVersion string
	GitCommit  string
	BuildDate  string
	GoVersion  string
	Compiler   string
	Platform   string
}

var Get = func() *Info {
	return &Info{
		Version:   version,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  platform,
	}
}

func ParseVersion(version string) (semver.Version, error) {
	version = strings.TrimSpace(version)
	v, err := semver.Parse(strings.TrimLeft(version, "v"))
	if err != nil {
		return semver.Version{}, fmt.Errorf("parsing semver: %w", err)
	}
	return v, nil
}
