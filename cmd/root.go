package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "siteintel",
	Short: "Deterministic local business website audit CLI",
	Long:  "siteintel audits public homepage HTML with deterministic extraction, scoring, and compact reports.",
}

// Execute runs the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
