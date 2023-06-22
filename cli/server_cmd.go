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
			AdminEmails:    []string{"libservice@ugent.be"},
			DeletedRecord:  "persistent",
			Granularity:    "YYYY-MM-DDThh:mm:ssZ",

			ListMetadataFormats: func(r *oaipmh.Request) ([]*oaipmh.MetadataFormat, error) {
				ctx := context.TODO()

				if r.Identifier != "" {
					formats, err := repo.GetRecordMetadataFormats(ctx, r.Identifier)
					if err == repository.ErrNotFound {
						return nil, oaipmh.ErrIDDoesNotExist
					}
					return formats, err
				}

				return repo.GetMetadataFormats(ctx)
			},

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
				if err == repository.ErrNotFound {
					return nil, oaipmh.ErrCannotDisseminateFormat
				}
				if err != nil {
					return nil, err
				}

				return rec, nil
			},

			ListRecords: func(r *oaipmh.Request) ([]*oaipmh.Record, *oaipmh.ResumptionToken, error) {
				ctx := context.TODO()

				exists, err := repo.HasMetadataFormat(ctx, r.MetadataPrefix)
				if err != nil {
					return nil, nil, err
				}
				if !exists {
					return nil, nil, oaipmh.ErrCannotDisseminateFormat
				}

				// if r.ResumptionToken != "" {
				// }

				recs, err := repo.GetRecords(ctx, r.MetadataPrefix)
				if err != nil {
					return nil, nil, err
				}

				return recs, nil, nil
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
