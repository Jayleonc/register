package cmd

import (
	"fmt"
	"github.com/Jayleonc/register/cmd/command"
	"github.com/spf13/cobra"
	"os"
)

const (
	cliName = "web"
)

var (
	rootCmd = &cobra.Command{
		Use: cliName,
	}
)

func init() {
	rootCmd.AddCommand(command.NewWebCommand())
}

func start() error {
	return rootCmd.Execute()
}

func MustStart() {
	if err := start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
