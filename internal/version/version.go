package version

// Build information. Populated at build-time via -ldflags
var (
	Version   = "unknown" // Version
	Commit    = "unknown" // Git SHA
	BuildTime = "unknown" // Build timestamp
)

// Info holds version information
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
}

// GetInfo returns all version information
func GetInfo() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
	}
}
func init() {
	println("Version:", Version)
	println("Commit:", Commit)
	println("BuildTime:", BuildTime)
}
