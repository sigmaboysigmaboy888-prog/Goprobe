package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/local/goprobe/internal/output"
	"github.com/local/goprobe/internal/runner"
)

func newTFTPCmd() *cobra.Command {
	cfgCopy := *globalCfg
	cfg := &cfgCopy

	var (
		host     string
		port     int
		wordlist string
	)

	cmd := &cobra.Command{
		Use:   "tftp",
		Short: "Enumerate files on a TFTP server",
		Long: `tftp – Brute-force file names on a TFTP server (UDP, RFC 1350).

Sends a Read Request (RRQ) for each word in the wordlist. A Data packet
response indicates the file exists; an Error packet or timeout means it does not.

Examples:
  goprobe tftp -H 192.168.1.1 -w filenames.txt
  goprobe tftp -H 192.168.1.1 -w filenames.txt --port 6969 --threads 20`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if host == "" {
				return fmt.Errorf("--host / -H is required")
			}
			if wordlist == "" {
				return fmt.Errorf("--wordlist / -w is required")
			}
			if err := applyGlobalConfig(cfg); err != nil {
				return err
			}

			cfg.Target = host
			cfg.TFTPHost = host
			cfg.TFTPPort = port
			cfg.Wordlist = wordlist

			output.Banner("tftp", fmt.Sprintf("%s:%d", host, port), wordlist, cfg.Threads)

			r := runner.NewTFTPRunner(cfg)
			_, err := r.Run(context.Background())
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&host, "host", "H", "", "TFTP server hostname or IP (required)")
	f.IntVar(&port, "port", 69, "TFTP UDP port")
	f.StringVarP(&wordlist, "wordlist", "w", "", "path to wordlist file (required)")

	return cmd
}
