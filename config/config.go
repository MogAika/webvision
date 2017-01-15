package config

type WebConfig struct {
	Host     string
	CertFile string
	KeyFile  string
}

type DBConfig struct {
	Dialect string
	Params  interface{}
}

type Config struct {
	DataPath    string
	FFmpeg      string
	DB          DBConfig
	Web         WebConfig
	Secret      string
	MaxDataSize int64
}
