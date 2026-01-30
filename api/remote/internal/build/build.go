package build

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// Info returns build metadata in a stable structure.
func Info() map[string]string {
	return map[string]string{
		"version":    Version,
		"commit":     Commit,
		"build_date": BuildDate,
	}
}
