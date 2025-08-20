package config

type ModulesConfig struct {
	ShellPool ShellPoolConfig
}

type ShellPoolConfig struct {
	Path string   // Shell binary path
	Args []string // Process arguments
	Size uint     // The maxmimum amount of shell processes to prepare
}
