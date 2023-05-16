package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/echox"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/versionx"
	"go.uber.org/zap"

	"go.infratographer.com/node-resolver/internal/config"
	"go.infratographer.com/node-resolver/internal/graphapi"
)

var defaultLBAPIListenAddr = ":8080"

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the location API",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	echox.MustViperFlags(viper.GetViper(), serveCmd.Flags(), defaultLBAPIListenAddr)
}

func serve(ctx context.Context) {
	err := otelx.InitTracer(config.AppConfig.Tracing, appName, logger)
	if err != nil {
		logger.Fatalw("failed to initialize tracer", "error", err)
	}

	srv, err := echox.NewServer(
		logger.Desugar(),
		echox.Config{
			Listen:              viper.GetString("server.listen"),
			ShutdownGracePeriod: viper.GetDuration("server.shutdown-grace-period"),
		},
		versionx.BuildDetails(),
	)
	if err != nil {
		logger.Fatalw("failed to create server", zap.Error(err))
	}

	schema, err := os.ReadFile("schema.graphql")
	if err != nil {
		logger.Fatalw("failed to read graphql schema file", "error", err)
	}

	r, err := graphapi.NewResolver(logger.Named("resolvers"), string(schema))
	if err != nil {
		logger.Fatalw("failed to create graphql resolver", "error", err)
	}

	srv.AddHandler(r)

	if err := srv.RunWithContext(ctx); err != nil {
		logger.Errorw("failed to run server", "error", zap.Error(err))
	}
}
