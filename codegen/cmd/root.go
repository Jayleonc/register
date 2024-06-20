package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codegen",
	Short: "Code generation tool for service clients",
	Long:  `A tool to generate client code for services registered in the etcd registry.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
