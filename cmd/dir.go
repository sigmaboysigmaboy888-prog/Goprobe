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

func newDirCmd() *cobra.Command {
	// Dereference into a local copy so this subcommand's flags
	// are isolated from the shared pointer.
	cfgCopy := *globalCfg
	cfg := &cfgCopy

	var (
		target     string
		wordlist   string
		extensions string
		addSlash   bool
		method     string
		postData   string
	)

	cmd := &cobra.Command{
		Use:   "dir",
		Short: "Brute-force directories and files",
		Long: `dir – Brute-force hidden directories and files on an HTTP/S server.

Each word in the wordlist is appended to the base URL. Optional extensions
(-x) generate additional requests per word (e.g. word, word.php, word.html).
Use --add-slash to also try a trailing slash variant.

Examples:
  goprobe dir -u https://target.com -w common.txt
  goprobe dir -u https://target.com -w common.txt -x php,html,txt
  goprobe dir -u https://target.com -w common.txt --status-codes 200,301,302
  goprobe dir -u https://target.com -w common.txt -H "X-Forwarded-For: 127.0.0.1" --add-slash`,
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
			cfg.AddSlash = addSlash
			cfg.Method = strings.ToUpper(method)
			cfg.PostData = postData

			if extensions != "" {
				cfg.Extensions = strings.Split(extensions, ",")
			}

			output.Banner("dir", target, wordlist, cfg.Threads)

			r := runner.NewDirRunner(cfg)
			_, err := r.Run(context.Background())
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&target, "url", "u", "", "target URL (required)")
	f.StringVarP(&wordlist, "wordlist", "w", "", "path to wordlist file (required)")
	f.StringVarP(&extensions, "extensions", "x", "", "file extensions to append, comma-separated (e.g. php,html,txt)")
	f.BoolVar(&addSlash, "add-slash", false, "also probe each word with a trailing slash")
	f.StringVarP(&method, "method", "m", "GET", "HTTP method (GET, POST, HEAD, …)")
	f.StringVar(&postData, "data", "", "POST body data")

	return cmd
}
