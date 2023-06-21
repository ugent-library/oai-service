package cli

import (
	"context"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	"github.com/ugent-library/httpx/render"
	"github.com/ugent-library/oai-service/models"
	"github.com/ugent-library/oai-service/oaipmh"
	"github.com/ugent-library/oai-service/repository"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// setup services
		repo, err := repository.New(config.Repo.Conn)
		if err != nil {
			return err
		}

		// setup oai provider
		// TODO simpler verb request and response types?
		oaiProvider, err := oaipmh.NewProvider(oaipmh.ProviderConfig{
			ErrorHandler:   func(err error) { logger.Error(err) },
			RepositoryName: "Ghent University Institutional Archive",
			BaseURL:        "https://biblio.ugent.be/oai",
			AdminEmail:     []string{"libservice@ugent.be"},
			DeletedRecord:  "persistent",
			Granularity:    "YYYY-MM-DDThh:mm:ssZ",
			GetRecord: func(r *oaipmh.Request) (*oaipmh.Record, error) {
				ctx := context.TODO()

				exists, err := repo.HasRecord(ctx, r.Identifier)
				if err != nil {
					return nil, err
				}
				if !exists {
					return nil, oaipmh.ErrIDDoesNotExist
				}

				rec, err := repo.GetRecord(ctx, r.Identifier, r.MetadataPrefix)
				if err == models.ErrNotFound {
					return nil, oaipmh.ErrCannotDisseminateFormat
				}
				if err != nil {
					return nil, err
				}

				if rec.Deleted {
					return &oaipmh.Record{
						Header: &oaipmh.Header{
							Identifier: rec.Identifier,
							Datestamp:  rec.Datestamp.UTC().Format(time.RFC3339),
							Status:     "deleted",
							SetSpec:    rec.SetSpecs,
						},
					}, nil
				}

				return &oaipmh.Record{
					Header: &oaipmh.Header{
						Identifier: rec.Identifier,
						Datestamp:  rec.Datestamp.UTC().Format(time.RFC3339),
						SetSpec:    rec.SetSpecs,
					},
					Metadata: &oaipmh.Payload{
						XML: rec.Metadata,
					},
				}, nil
			},
		})
		if err != nil {
			return err
		}

		// setup health checker
		// TODO add checkers
		healthChecker := health.NewChecker()

		// setup router
		router := chi.NewMux()
		router.Use(middleware.RequestID)
		if config.Env != "local" {
			router.Use(middleware.RealIP)
		}
		router.Use(zaphttp.SetLogger(logger.Desugar(), zapchi.RequestID))
		router.Use(middleware.RequestLogger(zapchi.LogFormatter()))
		router.Use(middleware.Recoverer)
		router.Use(middleware.StripSlashes)

		router.Get("/health", health.NewHandler(healthChecker))
		router.Get("/info", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, http.StatusOK, &struct {
				Branch string `json:"branch,omitempty"`
				Commit string `json:"commit,omitempty"`
			}{
				Branch: config.Source.Branch,
				Commit: config.Source.Commit,
			})
		})
		router.Method("GET", "/", oaiProvider)

		// start server
		server := graceful.WithDefaults(&http.Server{
			Addr:         config.Addr(),
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		})
		logger.Infof("starting server at %s", config.Addr())
		if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
			return err
		}
		logger.Info("gracefully stopped server")

		return nil
	},
}
