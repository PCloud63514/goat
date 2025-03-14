package profile

var (
	std = New()
)

func Get() []string {
	return std.Get()
}

func Contains(profile string) bool {
	return std.Contains(profile)
}
