package profile

import (
	"flag"
	"os"
	"strings"
)

const (
	ProfileDefault   = "default"
	ProfileFlag      = "profile"
	ProfileShortFlag = "p"
	ProfileSep       = ","
)

type Profile struct {
	value []string
}

func New() *Profile {
	profile := mergeSlicesUnique([]string{ProfileDefault}, readProfiles())
	return &Profile{
		value: profile,
	}
}

func (p *Profile) Get() []string {
	return p.value
}

func (p *Profile) Contains(profile string) bool {
	for _, v := range p.value {
		if v == profile {
			return true
		}
	}
	return false
}

func readProfiles() []string {
	cfs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	p1 := cfs.String(ProfileFlag, "", "Comma-separated list of profiles")
	p2 := cfs.String(ProfileShortFlag, "", "Comma-separated list of profiles")
	cfs.SetOutput(&strings.Builder{})

	if err := cfs.Parse(os.Args[1:]); err != nil {
		return make([]string, 0)
	}

	if p1 != nil && *p1 != "" {
		_profiles := strings.Replace(*p1, ProfileDefault, "", 0)
		return strings.Split(_profiles, ",")
	}

	if p2 != nil && *p2 != "" {
		_p := strings.Replace(*p2, ProfileDefault, "", 0)
		return strings.Split(_p, ProfileSep)
	}

	return make([]string, 0)
}

func mergeSlicesUnique(a, b []string) []string {
	m := make(map[string]bool)
	var result []string

	for _, item := range a {
		if _, ok := m[item]; !ok {
			m[item] = true
			result = append(result, item)
		}
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			m[item] = true
			result = append(result, item)
		}
	}

	return result
}
