package project

import "github.com/spf13/cobra"

var projectMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate a project",
}

func init() {
	projectRootCmd.AddCommand(projectMigrateCmd)
}
