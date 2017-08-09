package cfg

type Configuration struct {
	Verbose bool
	Files   []string
	Command string
}

var (
	Global Configuration
)
