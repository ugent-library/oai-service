package cli

import "fmt"

type Config struct {
	// Version info
	Source struct {
		Branch string `env:"BRANCH"`
		Commit string `env:"COMMIT"`
	} `envPrefix:"SOURCE_"`
	// Env must be local, development, test or production
	Env  string `env:"OAI_ENV" envDefault:"production"`
	Host string `env:"OAI_HOST"`
	Port int    `env:"OAI_PORT" envDefault:"3000"`
	Repo struct {
		Conn   string `env:"CONN,notEmpty"`
		Secret string `env:"SECRET,notEmpty"`
	} `envPrefix:"OAI_REPO_"`
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
