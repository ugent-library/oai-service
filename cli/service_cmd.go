package cli

import (
	"runtime"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/spf13/cobra"
	"github.com/ugent-library/oai-service/repositories"
	"github.com/ugent-library/oai-service/services/oai"
)

func init() {
	rootCmd.AddCommand(serviceCmd)
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Start NATS service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := startService(); err != nil {
			return err
		}
		runtime.Goexit()
		return nil
	},
}

func startService() error {
	// setup repo
	repo, err := repositories.New(repositories.Config{
		Conn:   config.Repo.Conn,
		Secret: []byte(config.Repo.Secret),
	})
	if err != nil {
		return err
	}

	// setup nats connection
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return err
	}

	oaiService := oai.NewService(repo)

	srv, err := micro.AddService(nc, oai.Config)
	if err != nil {
		return err
	}

	grp := srv.AddGroup("oai")

	if err := grp.AddEndpoint("AddMetadataFormat", micro.HandlerFunc(oaiService.AddMetadataFormat)); err != nil {
		return err
	}
	if err := grp.AddEndpoint("AddSet", micro.HandlerFunc(oaiService.AddSet)); err != nil {
		return err
	}
	if err := grp.AddEndpoint("AddRecord", micro.HandlerFunc(oaiService.AddRecord)); err != nil {
		return err
	}
	if err := grp.AddEndpoint("DeleteRecord", micro.HandlerFunc(oaiService.DeleteRecord)); err != nil {
		return err
	}

	return nil
}
