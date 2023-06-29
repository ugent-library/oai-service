package cli

import (
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/bufbuild/connect-go"

	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	"github.com/ugent-library/httpx/render"
	"github.com/ugent-library/oai-service/gen/oai/v1/oaiv1connect"
	"github.com/ugent-library/oai-service/grpcserver"
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
			Sets:           true,
			ErrorHandler:   func(err error) { logger.Error(err) },
			Backend:        repo,
		})
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
			}{
				Branch: config.Source.Branch,
				Commit: config.Source.Commit,
			})
		})

		// mount oai provider
		mux.Get("/oai.xsl", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "public/oai.xsl")
		})
		mux.Method("GET", "/", oaiProvider)

		// mount grpc server
		grpcReflector := grpcreflect.NewStaticReflector(oaiv1connect.OaiServiceName)
		grpcChecker := grpchealth.NewStaticChecker(oaiv1connect.OaiServiceName)
		mux.Mount(oaiv1connect.NewOaiServiceHandler(
			grpcserver.NewServer(repo),
			connect.WithInterceptors(
				grpcserver.NewAuthInterceptor(grpcserver.AuthConfig{
					Token: config.GRPC.Secret,
				}),
			),
		))
		mux.Mount(grpcreflect.NewHandlerV1(grpcReflector))
		mux.Mount(grpcreflect.NewHandlerV1Alpha(grpcReflector))
		mux.Mount(grpchealth.NewHandler(grpcChecker))

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
