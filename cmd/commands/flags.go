package commands

const (
	// FlagPort command flag: port for listener
	FlagPort  = "port"

	// FlagHost command flag: remote hostname
	FlagHost  = "host"

	// FlagTo command flag: path to copy files on remote
	FlagTo = "to"

	// FlagFrom command flag: path to copy target (file or directory)
	FlagFrom = "from"

	// FlagConfig command flag: path to serve configs
	FlagConfig = "config"

	// FlagLimit global flag: send goroutines limit
	FlagLimit = "limit"
)