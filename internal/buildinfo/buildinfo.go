package buildinfo

// Эти переменные можно прокидывать через -ldflags на CI/CD
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
