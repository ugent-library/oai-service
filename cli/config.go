package cli

import "fmt"

// Version info
type Version struct {
	Branch string `env:"SOURCE_BRANCH"`
	Commit string `env:"SOURCE_COMMIT"`
	Image  string `env:"IMAGE_NAME"`
}

type Config struct {
	// Env must be local, development, test or production
	Env    string `env:"ENV" envDefault:"production"`
	Host   string `env:"HOST"`
	Port   int    `env:"PORT" envDefault:"3000"`
	APIKey string `env:"API_KEY"`
	Repo   struct {
		Conn   string `env:"CONN,notEmpty"`
		Secret string `env:"SECRET,notEmpty"`
	} `envPrefix:"REPO_"`
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
