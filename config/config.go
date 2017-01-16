package config

type Web struct {
	Host     string
	CertFile string
	KeyFile  string
}

type DB struct {
	Dialect string
	Params  interface{}
}

type Cookie struct {
	HashKey  string
	BlockKey string
	LifeTime int64
}

type Config struct {
	DataPath    string
	FFmpeg      string
	DB          DB
	Web         Web
	Cookie      Cookie
	MaxDataSize int64
	Secret      string
}
