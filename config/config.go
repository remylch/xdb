package config

type SecretConfig struct {
	HashKey string
}

type ServerConfig struct {
	NodeAddr       string
	ApiAddr        string
	BootstrapNodes []string
}

type StorageConfig struct {
	DataDir string
}

type LogConfig struct {
	LogDir string
}

type GlobalConfig struct {
	Secret  SecretConfig
	Server  ServerConfig
	Storage StorageConfig
	Log     LogConfig
}
