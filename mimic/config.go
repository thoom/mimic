package mimic

type Config struct {
	Home           string `env:"HOME"`
	HostURL        string `env:"HOST_URL" envDefault:"http://localhost:8778"`
	Password       string `env:"PASSWORD,unset"`
	ParsedHostName struct {
		Protocol string
		Host     string
		Port     string
	}
}
