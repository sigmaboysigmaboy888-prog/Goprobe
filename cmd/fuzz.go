package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/local/goprobe/internal/output"
	"github.com/local/goprobe/internal/runner"
)

func newFuzzCmd() *cobra.Command {
	cfgCopy := *globalCfg
	cfg := &cfgCopy

	var (
		target      string
		wordlist    string
		method      string
		postData    string
		fuzzKeyword string
	)

	cmd := &cobra.Command{
		Use:   "fuzz",
		Short: "Generic HTTP fuzzer with a FUZZ keyword placeholder",
		Long: `fuzz – Replace a keyword token anywhere in the request with each wordlist entry.

Place the keyword (default: FUZZ) in the URL, headers (--header), cookies
(--cookies), or POST body (--data). goprobe will substitute the keyword and
send a request for each word.

Examples:
  # Path fuzzing
  goprobe fuzz -u https://example.com/FUZZ -w wordlist.txt

  # Query string fuzzing
  goprobe fuzz -u "https://example.com/api?id=FUZZ" -w ids.txt

  # POST body fuzzing
  goprobe fuzz -u https://example.com/login -m POST \
    --data "username=admin&password=FUZZ" -w passwords.txt

  # Header fuzzing with custom keyword
  goprobe fuzz -u https://example.com/ -w tokens.txt \
    -H "Authorization: Bearer FUZZ" --fuzz-keyword FUZZ

  # Filter to only show specific status codes
  goprobe fuzz -u https://example.com/FUZZ -w words.txt --status-codes 200,302`,
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
			cfg.Method = strings.ToUpper(method)
			cfg.PostData = postData
			cfg.FuzzKeyword = fuzzKeyword

			output.Banner("fuzz", target, wordlist, cfg.Threads)

			r := runner.NewFuzzRunner(cfg)
			_, err := r.Run(context.Background())
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&target, "url", "u", "", "target URL containing the FUZZ keyword (required)")
	f.StringVarP(&wordlist, "wordlist", "w", "", "path to wordlist file (required)")
	f.StringVarP(&method, "method", "m", "GET", "HTTP method (GET, POST, PUT, …)")
	f.StringVar(&postData, "data", "", "POST/PUT body; place FUZZ keyword here to fuzz body")
	f.StringVar(&fuzzKeyword, "fuzz-keyword", "FUZZ", "placeholder token to substitute in the request")

	return cmd
}
