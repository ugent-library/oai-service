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
	"github.com/ugent-library/oai-service/api"
	"github.com/ugent-library/oai-service/gen/oai/v1/oaiv1connect"
	"github.com/ugent-library/oai-service/oaipmh"
	"github.com/ugent-library/oai-service/repository"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
			StyleSheet:     "/oai.xsl",

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

				if r.ResumptionToken != "" {
					recs, token, err := repo.GetMoreRecords(ctx, r.ResumptionToken)
					if err != nil {
						return nil, nil, err
					}
					return recs, token, nil
				}

				exists, err := repo.HasMetadataFormat(ctx, r.MetadataPrefix)
				if err != nil {
					return nil, nil, err
				}
				if !exists {
					return nil, nil, oaipmh.ErrCannotDisseminateFormat
				}

				recs, token, err := repo.GetRecords(ctx, r.MetadataPrefix, r.Set, r.From, r.Until)
				if err != nil {
					return nil, nil, err
				}
				return recs, token, nil
			},
		})
		if err != nil {
			return err
		}

		// setup api server
		apiPath, apiHandler := oaiv1connect.NewOaiServiceHandler(api.NewServer(repo))

		// setup health checker
		// TODO add checkers
		healthChecker := health.NewChecker()

		// setup mux
		mux := chi.NewMux()
		mux.Use(middleware.RequestID)
		if config.Env != "local" {
			mux.Use(middleware.RealIP)
		}
		mux.Use(zaphttp.SetLogger(logger.Desugar(), zapchi.RequestID))
		mux.Use(middleware.RequestLogger(zapchi.LogFormatter()))
		mux.Use(middleware.Recoverer)

		mux.Get("/health", health.NewHandler(healthChecker))
		mux.Get("/info", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, http.StatusOK, &struct {
				Branch string `json:"branch,omitempty"`
				Commit string `json:"commit,omitempty"`
			}{
				Branch: config.Source.Branch,
				Commit: config.Source.Commit,
			})
		})
		mux.Get("/oai.xsl", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "public/oai.xsl")
		})
		mux.Method("GET", "/", oaiProvider)
		mux.Mount(apiPath, apiHandler)

		handler := h2c.NewHandler(mux, &http2.Server{})

		// start server
		server := graceful.WithDefaults(&http.Server{
			Addr:         config.Addr(),
			Handler:      handler,
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
