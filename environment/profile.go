package environment

import (
	"flag"
	"os"
	"strings"
)

const (
	profileDefault   = "default"
	profileSystem    = "system"
	profileFlag      = "profile"
	profileShortFlag = "p"
	profileSep       = ","
)

var (
	defaultProfiles = []string{profileDefault, profileSystem}
)

func readProfiles() []string {
	cfs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	p1 := cfs.String(profileFlag, "", "Comma-separated list of profiles")
	p2 := cfs.String(profileShortFlag, "", "Comma-separated list of profiles")
	cfs.SetOutput(&strings.Builder{})

	if err := cfs.Parse(os.Args[1:]); err != nil {
		return make([]string, 0)
	}

	if p1 != nil && *p1 != "" {
		_profiles := strings.Replace(*p1, profileDefault, "", 0)
		_profiles = strings.Replace(*p1, profileSystem, "", 0)
		return strings.Split(_profiles, ",")
	}

	if p2 != nil && *p2 != "" {
		_p := strings.Replace(*p2, profileDefault, "", 0)
		_p = strings.Replace(*p2, profileSystem, "", 0)
		return strings.Split(_p, profileSep)
	}

	return make([]string, 0)
}
