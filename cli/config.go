package cli

import "fmt"

type Config struct {
	// Version info
	Version struct {
		Branch string `env:"SOURCE_BRANCH"`
		Commit string `env:"SOURCE_COMMIT"`
		Image  string `env:"IMAGE_NAME"`
	}
	// Env must be local, development, test or production
	Env  string `env:"OAI_ENV" envDefault:"production"`
	Host string `env:"OAI_HOST"`
	Port int    `env:"OAI_PORT" envDefault:"3000"`
	Repo struct {
		Conn   string `env:"CONN,notEmpty"`
		Secret string `env:"SECRET,notEmpty"`
	} `envPrefix:"OAI_REPO_"`
	GRPC struct {
		Secret string `env:"SECRET,notEmpty"`
	} `envPrefix:"OAI_GRPC_"`
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
