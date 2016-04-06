package settings

type WebSettings struct {
	Host        string // 0.0.0.0:80
	TlsHost     string // 0.0.0.0:443
	TlsCertFile string
	TlsKeyFile  string
}

type DBSettings struct {
	Dialect string
	Params  interface{}
}

type Settings struct {
	DataPath    string
	FFmpeg      string
	DB          DBSettings
	Web         WebSettings
	Secret      string
	MaxDataSize int64
}
