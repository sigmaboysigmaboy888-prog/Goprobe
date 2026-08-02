package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/local/goprobe/internal/output"
	"github.com/local/goprobe/internal/runner"
)

func newDNSCmd() *cobra.Command {
	cfgCopy := *globalCfg
	cfg := &cfgCopy

	var (
		domain   string
		wordlist string
		resolver string
	)

	cmd := &cobra.Command{
		Use:   "dns",
		Short: "Enumerate subdomains via DNS resolution",
		Long: `dns – Brute-force subdomains by resolving <word>.<domain> for each wordlist entry.

A wildcard check is performed before scanning begins. If the target domain
resolves wildcard requests, a warning is shown and results may contain false
positives. Use a custom resolver with --resolver to bypass local DNS caches.

Examples:
  goprobe dns -d example.com -w subdomains.txt
  goprobe dns -d example.com -w subdomains.txt --threads 100
  goprobe dns -d example.com -w subdomains.txt --resolver 8.8.8.8`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if domain == "" {
				return fmt.Errorf("--domain / -d is required")
			}
			if wordlist == "" {
				return fmt.Errorf("--wordlist / -w is required")
			}
			if err := applyGlobalConfig(cfg); err != nil {
				return err
			}

			cfg.Target = domain
			cfg.Wordlist = wordlist
			cfg.Resolver = resolver

			output.Banner("dns", domain, wordlist, cfg.Threads)

			r := runner.NewDNSRunner(cfg)
			_, err := r.Run(context.Background())
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&domain, "domain", "d", "", "target domain to enumerate (required)")
	f.StringVarP(&wordlist, "wordlist", "w", "", "path to wordlist file (required)")
	f.StringVar(&resolver, "resolver", "", "custom DNS resolver IP[:port] (e.g. 8.8.8.8 or 1.1.1.1:53)")

	return cmd
}
