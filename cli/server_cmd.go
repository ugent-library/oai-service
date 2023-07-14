package cli

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	"github.com/ugent-library/httpx/render"
	"github.com/ugent-library/oai-service/api/v1"
	"github.com/ugent-library/oai-service/oaipmh"
	"github.com/ugent-library/oai-service/repositories"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

type apiSecurityHandler struct{}

func (s *apiSecurityHandler) HandleApiKey(ctx context.Context, operationName string, t api.ApiKey) (context.Context, error) {
	if t.APIKey == "123" {
		return ctx, nil
	}
	return ctx, errors.New("go away")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// setup repo
		repo, err := repositories.New(repositories.Config{
			Conn:   config.Repo.Conn,
			Secret: []byte(config.Repo.Secret),
		})
		if err != nil {
			return err
		}

		// setup oai provider
		oaiProvider, err := oaipmh.NewProvider(oaipmh.ProviderConfig{
			RepositoryName: "Ghent University Institutional Archive",
			BaseURL:        "https://biblio.ugent.be/oai",
			AdminEmails:    []string{"libservice@ugent.be"},
			DeletedRecord:  "persistent",
			Granularity:    "YYYY-MM-DDThh:mm:ssZ",
			StyleSheet:     "/oai.xsl",
			ErrorHandler:   func(err error) { logger.Error(err) },
			Backend:        repo,
		})
		if err != nil {
			return err
		}

		// setup api
		apiServer, err := api.NewServer(api.NewService(repo), &apiSecurityHandler{})
		if err != nil {
			return err
		}

		// setup mux
		mux := chi.NewMux()
		mux.Use(middleware.RequestID)
		if config.Env != "local" {
			mux.Use(middleware.RealIP)
		}
		mux.Use(zaphttp.SetLogger(logger.Desugar(), zapchi.RequestID))
		mux.Use(middleware.RequestLogger(zapchi.LogFormatter()))
		mux.Use(middleware.Recoverer)

		// mount health and info
		mux.Get("/health", health.NewHandler(health.NewChecker())) // TODO add checkers
		mux.Get("/info", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, http.StatusOK, &struct {
				Branch string `json:"branch,omitempty"`
				Commit string `json:"commit,omitempty"`
				Image  string `json:"image,omitempty"`
			}{
				Branch: config.Version.Branch,
				Commit: config.Version.Commit,
				Image:  config.Version.Image,
			})
		})

		// mount api
		mux.Mount("/api/v1", http.StripPrefix("/api/v1", apiServer))
		mux.Get("/api/v1/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "api/v1/openapi.yaml")
		})
		mux.Mount("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("public/swagger-ui-5.1.0"))))

		// mount oai provider
		mux.Method("GET", "/", oaiProvider)
		mux.Get("/oai.xsl", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "public/oai.xsl")
		})

		// start server
		server := graceful.WithDefaults(&http.Server{
			Addr:         config.Addr(),
			Handler:      h2c.NewHandler(mux, &http2.Server{}),
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
