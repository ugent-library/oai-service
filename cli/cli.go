package cli

import (
	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	version Version
	config  Config
	logger  *zap.SugaredLogger

	rootCmd = &cobra.Command{
		Use:   "oai",
		Short: "OAI CLI",
	}
)

func init() {
	cobra.OnInitialize(initVersion, initConfig, initLogger)
	cobra.OnFinalize(func() {
		logger.Sync()
	})
}

func initVersion() {
	cobra.CheckErr(env.Parse(&version))
}

func initConfig() {
	cobra.CheckErr(env.ParseWithOptions(&config, env.Options{
		Prefix: "OAI_",
	}))
}

func initLogger() {
	if config.Env == "local" {
		l, err := zap.NewDevelopment()
		cobra.CheckErr(err)
		logger = l.Sugar()
	} else {
		l, err := zap.NewProduction()
		cobra.CheckErr(err)
		logger = l.Sugar()
	}
}

func Run() {
	cobra.CheckErr(rootCmd.Execute())
}
