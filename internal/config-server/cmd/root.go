package cmd

import (
	"gitlab.snapp.ir/snappcloud/config-server/internal/config"
	"gitlab.snapp.ir/snappcloud/config-server/internal/config-server/cmd/api"

	"github.com/spf13/cobra"
)

// NewRootCommand creates a new config-server root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "config-server",
	}

	cfg := config.New()

	api.Register(root, cfg)

	return root
}
