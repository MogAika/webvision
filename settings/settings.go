package settings

type WebSettings struct {
	Host        string // 0.0.0.0:8080
	Url         string // https://www.url.com
	Tls         bool
	TlsCertFile string
	TlsKeyFile  string
}

type DBSettings struct {
	Dialect string
	Params  interface{}
}

type Settings struct {
	DataPath    string
	DB          DBSettings
	Web         WebSettings
	Secret      string
	MaxDataSize int64
}
