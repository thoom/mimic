package mimic

import "time"

type Config struct {
	Home         string         `env:"HOME"`
	Host         string         `env:"HOST" envDefault:"localhost"`
	Port         int            `env:"PORT" envDefault:"8778"`
	Password     string         `env:"PASSWORD,unset"`
	IsProduction bool           `env:"PRODUCTION"`
	Duration     time.Duration  `env:"DURATION"`
	Hosts        []string       `env:"HOSTS" envSeparator:":"`
	TempFolder   string         `env:"TEMP_FOLDER,expand" envDefault:"${HOME}/tmp"`
	StringInts   map[string]int `env:"MAP_STRING_INT"`
}
