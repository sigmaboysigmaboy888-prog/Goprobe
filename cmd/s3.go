package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/local/goprobe/internal/output"
	"github.com/local/goprobe/internal/runner"
)

func newS3Cmd() *cobra.Command {
	cfgCopy := *globalCfg
	cfg := &cfgCopy

	var (
		wordlist string
		region   string
		scheme   string
	)

	cmd := &cobra.Command{
		Use:   "s3",
		Short: "Discover exposed AWS S3 buckets",
		Long: `s3 – Brute-force AWS S3 bucket names using the public REST endpoint.

Each word in the wordlist is tested as a bucket name. Buckets that respond
with 200 (public) or 403 (exists but private) are reported.

Set --region to a specific AWS region string, or leave it as the default
(s3.amazonaws.com) for us-east-1 style endpoints.

Examples:
  goprobe s3 -w buckets.txt
  goprobe s3 -w buckets.txt --region us-west-2
  goprobe s3 -w buckets.txt --scheme http`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wordlist == "" {
				return fmt.Errorf("--wordlist / -w is required")
			}
			if err := applyGlobalConfig(cfg); err != nil {
				return err
			}

			cfg.Target = region
			cfg.Wordlist = wordlist
			cfg.Scheme = scheme

			output.Banner("s3", "amazonaws.com/"+region, wordlist, cfg.Threads)

			r := runner.NewS3Runner(cfg)
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
	f.StringVar(&region, "region", "s3.amazonaws.com",
		"AWS region string or custom S3 endpoint (e.g. us-east-1, us-west-2)")
	f.StringVar(&scheme, "scheme", "https", "URL scheme: https or http")

	return cmd
}
