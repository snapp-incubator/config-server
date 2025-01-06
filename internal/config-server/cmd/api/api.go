package api

import (
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gitlab.snapp.ir/snappcloud/config-server/internal/engine"
	"gitlab.snapp.ir/snappcloud/config-server/internal/handler"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"gitlab.snapp.ir/snappcloud/config-server/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main(cfg config.Config) {
	clusterConfig, err := getClusterConfig()
	if err != nil {
		log.Errorf("Failed loading kubeconfig: %s", err)
		os.Exit(1)
	}

	configEngine := engine.NewEngine(clusterConfig)
	configHandler := handler.NewConfigHandler(*configEngine, cfg.AppConfigs)

	app := echo.New()

	app.GET("/healthz", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	api := app.Group("/api")
	api.GET("/:endpoint", configHandler.Get)

	if err := app.Start(fmt.Sprintf(":%d", cfg.API.Port)); !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("echo initiation failed: %s", err)
	}

	logrus.Println("API has been started :D")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// Register API command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "api",
			Short: "Run API to serve the requests",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}

func getClusterConfig() (*rest.Config, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err == nil {
		return clusterConfig, nil
	}

	kubeconfig, err := kubeconfigLocation()
	if err != nil {
		return nil, err
	}

	clusterConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return clusterConfig, nil
}

func kubeconfigLocation() (string, error) {
	value, present := os.LookupEnv("KUBECONFIG")
	if present {
		fileExist, err := exists(value)
		if err != nil {
			return "", err
		}
		if fileExist {
			return value, nil
		}
	}
	return filepath.Join(os.Getenv("HOME"), ".kube", "config"), nil
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
