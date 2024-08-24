package config

type Config struct {
	Routing struct {
		Port string
	}
}

func NewConfig() *Config {
	c := new(Config)
	c.Routing.Port = ":8000"

	return c
}
