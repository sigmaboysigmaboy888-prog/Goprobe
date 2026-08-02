package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/local/goprobe/internal/output"
	"github.com/local/goprobe/internal/runner"
)

func newGCSCmd() *cobra.Command {
	cfgCopy := *globalCfg
	cfg := &cfgCopy

	var (
		wordlist string
		scheme   string
	)

	cmd := &cobra.Command{
		Use:   "gcs",
		Short: "Discover exposed Google Cloud Storage buckets",
		Long: `gcs – Brute-force Google Cloud Storage bucket names.

Each word is tested against both the path-style URL
(https://storage.googleapis.com/<bucket>) and the subdomain-style URL
(https://<bucket>.storage.googleapis.com).

Buckets responding with 200 (public) or 403 (exists but private) are reported.

Examples:
  goprobe gcs -w buckets.txt
  goprobe gcs -w buckets.txt --threads 30 --scheme http`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wordlist == "" {
				return fmt.Errorf("--wordlist / -w is required")
			}
			if err := applyGlobalConfig(cfg); err != nil {
				return err
			}

			cfg.Wordlist = wordlist
			cfg.Scheme = scheme

			output.Banner("gcs", "storage.googleapis.com", wordlist, cfg.Threads)

			r := runner.NewGCSRunner(cfg)
			_, err := r.Run(context.Background())
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&wordlist, "wordlist", "w", "", "path to wordlist file (required)")
	f.StringVar(&scheme, "scheme", "https", "URL scheme: https or http")

	return cmd
}
