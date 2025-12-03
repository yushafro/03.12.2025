package config

import "time"

type Config struct {
	Host               string
	Port               string
	ReadHTO            time.Duration `mapstructure:"read_header_timeout"`
	ReadTO             time.Duration `mapstructure:"read_timeout"`
	WriteTO            time.Duration `mapstructure:"write_timeout"`
	IdleTO             time.Duration `mapstructure:"idle_timeout"`
	ClientTO           time.Duration `mapstructure:"client_timeout"`
	FileRepositoryPath string        `mapstructure:"file_repository_path"`
}
