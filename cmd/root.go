package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/local/goprobe/internal/output"
	"github.com/local/goprobe/internal/runner"
)

const defaultUserAgent = "goprobe/1.0 (https://github.com/local/goprobe)"

// globalCfg is the shared Config populated by persistent root flags.
var globalCfg = &runner.Config{
	Threads:   10,
	Timeout:   10 * time.Second,
	Method:    "GET",
	UserAgent: defaultUserAgent,
}

var rootCmd = &cobra.Command{
	Use:   "goprobe",
	Short: "Fast web & DNS recon tool",
	Long: `goprobe – fast, multi-mode web and DNS brute-force scanner

Modes:
  dir    Brute-force directories and files on an HTTP server
  dns    Enumerate subdomains via DNS resolution
  vhost  Discover virtual hosts by fuzzing the Host header
  s3     Discover exposed AWS S3 buckets
  gcs    Discover exposed Google Cloud Storage buckets
  tftp   Enumerate files on a TFTP server
  fuzz   Generic HTTP fuzzer with a FUZZ keyword placeholder

Examples:
  goprobe dir  -u https://example.com -w /usr/share/wordlists/common.txt
  goprobe dns  -d example.com        -w subdomains.txt --threads 50
  goprobe vhost -u http://10.0.0.1  -w vhosts.txt --append-domain example.com
  goprobe s3   --region us-east-1   -w buckets.txt
  goprobe fuzz -u https://example.com/FUZZ -w payloads.txt -H "Authorization: Bearer TOKEN"`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute is the entry point called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func init() {
	pf := rootCmd.PersistentFlags()

	pf.IntVarP(&globalCfg.Threads, "threads", "t", 10,
		"number of concurrent workers")
	pf.DurationVar(&globalCfg.Timeout, "timeout", 10*time.Second,
		"per-request timeout (e.g. 5s, 500ms)")
	pf.BoolVarP(&globalCfg.Verbose, "verbose", "v", false,
		"show debug-level output (errors, skipped entries)")
	pf.BoolVarP(&globalCfg.Quiet, "quiet", "q", false,
		"suppress all output except found results")
	pf.StringVar(&globalCfg.StatusCodesRaw, "status-codes", "",
		"comma-separated HTTP status codes to report (default: anything except 404)")
	pf.BoolVar(&globalCfg.FollowRedirect, "follow-redirect", false,
		"follow HTTP redirects")
	pf.StringVar(&globalCfg.UserAgent, "user-agent", defaultUserAgent,
		"custom User-Agent header value")
	pf.StringArrayVarP(&globalCfg.Headers, "header", "H", nil,
		"custom header in \"Key: Value\" format (repeatable)")
	pf.StringVar(&globalCfg.Cookies, "cookies", "",
		"cookie string, e.g. \"session=abc123; token=xyz\"")
	pf.BoolVarP(&globalCfg.InsecureSkipVerify, "insecure", "k", false,
		"skip TLS certificate verification")

	// Register subcommands.
	rootCmd.AddCommand(
		newDirCmd(),
		newDNSCmd(),
		newVHostCmd(),
		newS3Cmd(),
		newGCSCmd(),
		newTFTPCmd(),
		newFuzzCmd(),
	)
}

// applyGlobalConfig applies verbose/quiet settings to the output package
// and parses raw status codes into the cfg struct.
func applyGlobalConfig(cfg *runner.Config) error {
	output.SetVerbose(cfg.Verbose)
	output.SetQuiet(cfg.Quiet)

	if cfg.StatusCodesRaw != "" {
		parts := strings.Split(cfg.StatusCodesRaw, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			code, err := strconv.Atoi(p)
			if err != nil {
				return fmt.Errorf("invalid status code %q: %w", p, err)
			}
			cfg.StatusCodes = append(cfg.StatusCodes, code)
		}
	}
	return nil
}
