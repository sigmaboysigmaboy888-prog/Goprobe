package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/local/goprobe/internal/output"
	"github.com/local/goprobe/internal/runner"
)

func newVHostCmd() *cobra.Command {
	cfgCopy := *globalCfg
	cfg := &cfgCopy

	var (
		target       string
		wordlist     string
		appendDomain string
	)

	cmd := &cobra.Command{
		Use:   "vhost",
		Short: "Discover virtual hosts by fuzzing the Host header",
		Long: `vhost – Send requests to a fixed server IP/URL while varying the Host header.

A baseline response size is measured first. Responses that differ from the
baseline are reported as potential virtual hosts. Use --append-domain to
automatically suffix each word with a domain name.

Examples:
  goprobe vhost -u http://10.0.0.1 -w words.txt
  goprobe vhost -u http://10.0.0.1 -w words.txt --append-domain internal.example.com
  goprobe vhost -u http://10.0.0.1 -w words.txt --status-codes 200,301`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if target == "" {
				return fmt.Errorf("--url / -u is required")
			}
			if wordlist == "" {
				return fmt.Errorf("--wordlist / -w is required")
			}
			if err := applyGlobalConfig(cfg); err != nil {
				return err
			}

			cfg.Target = target
			cfg.Wordlist = wordlist
			cfg.AppendDomain = appendDomain

			output.Banner("vhost", target, wordlist, cfg.Threads)

			r, err := runner.NewVHostRunner(cfg)
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}
			_, err = r.Run(context.Background())
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&target, "url", "u", "", "target URL/IP to probe (required)")
	f.StringVarP(&wordlist, "wordlist", "w", "", "path to wordlist file (required)")
	f.StringVar(&appendDomain, "append-domain", "", "suffix appended to each word, e.g. .example.com")

	return cmd
}
