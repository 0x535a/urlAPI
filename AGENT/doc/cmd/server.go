package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zhongxin/cmd/flags"
	"zhongxin/internal/conf"
	"zhongxin/internal/server"
	"zhongxin/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server",
	Long:  `Starts the server`,
	Run: func(cmd *cobra.Command, args []string) {
		serverStart()
	},
}

func serverStart() {
	Init()
	defer Release()
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	var r *gin.Engine
	if flags.Dev || flags.Beta {
		gin.SetMode(gin.DebugMode)
		r = gin.New()
		config.AllowAllOrigins = true
		config.AllowMethods = []string{"*"}
		r.Use(cors.New(config))
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}
	r.Use(
		server.RouterRecovery(),
		gin.LoggerWithWriter(log.StandardLogger().Out),
		gin.RecoveryWithWriter(log.StandardLogger().Out))

	server.Init(r)

	var httpHandler http.Handler = r
	httpBase := fmt.Sprintf("%s:%d", conf.Conf.Schema.Listen, conf.Conf.Schema.Port)
	util.Log.Infof("Starting HTTP server on %s", httpBase)
	httpServer := &http.Server{
		Addr:    httpBase,
		Handler: httpHandler,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			util.Log.Errorf("Failed to start server: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	util.Log.Infof("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		util.Log.Errorf("Failed to shutdown server: %v", err)
		return
	}
}

func init() {
	RootCmd.AddCommand(ServerCmd)
}
