package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	engine "siteintel/internal/audit"
	"siteintel/internal/report"
)

var auditFlags struct {
	niche     string
	format    string
	timeout   int
	maxBytes  int64
	userAgent string
	verbose   bool
}

var auditCmd = &cobra.Command{
	Use:   "audit <url>",
	Short: "Audit a public business homepage",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if auditFlags.format != "markdown" && auditFlags.format != "json" {
			return fmt.Errorf("unsupported format: %s", auditFlags.format)
		}
		result, err := engine.Run(context.Background(), args[0], engine.Options{
			Niche:          auditFlags.niche,
			TimeoutSeconds: auditFlags.timeout,
			MaxBytes:       auditFlags.maxBytes,
			UserAgent:      auditFlags.userAgent,
		})
		if err != nil {
			return err
		}
		if !auditFlags.verbose {
			result.Diagnostics.Warnings = nil
		}
		switch auditFlags.format {
		case "json":
			out, err := report.JSON(*result)
			if err != nil {
				return err
			}
			fmt.Fprint(os.Stdout, out)
		default:
			fmt.Fprint(os.Stdout, report.Markdown(*result))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
	auditCmd.Flags().StringVar(&auditFlags.niche, "niche", "local-service", "business niche: "+engine.SupportedNiches())
	auditCmd.Flags().StringVar(&auditFlags.format, "format", "markdown", "output format: markdown or json")
	auditCmd.Flags().IntVar(&auditFlags.timeout, "timeout", 15, "timeout in seconds for HTTP requests")
	auditCmd.Flags().Int64Var(&auditFlags.maxBytes, "max-bytes", 1500000, "maximum response size to process")
	auditCmd.Flags().StringVar(&auditFlags.userAgent, "user-agent", "siteintel/0.1 (+https://example.com)", "HTTP user agent")
	auditCmd.Flags().BoolVar(&auditFlags.verbose, "verbose", false, "show diagnostics")
}
