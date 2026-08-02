# goprobe

A fast, multi-mode web and DNS reconnaissance tool written in Go.

```
>> goprobe  fast web & dns recon
```

---

## Modes

| Mode  | What it does |
|-------|-------------|
| `dir`  | Brute-force hidden directories and files on an HTTP/S server |
| `dns`  | Enumerate subdomains via DNS resolution |
| `vhost` | Discover virtual hosts by fuzzing the `Host` header |
| `s3`   | Discover exposed AWS S3 buckets |
| `gcs`  | Discover exposed Google Cloud Storage buckets |
| `tftp` | Enumerate files on a TFTP server (UDP/RFC 1350) |
| `fuzz` | Generic HTTP fuzzer using a `FUZZ` keyword placeholder |

---

## Build

Requires Go 1.22 or later.

```bash
git clone https://github.com/local/goprobe
cd goprobe
go mod download
go build -o goprobe .
```

Cross-compile for Linux from macOS/Windows:

```bash
GOOS=linux GOARCH=amd64 go build -o goprobe-linux .
```

---

## Global Flags

These flags work across all modes:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--threads` | `-t` | `10` | Number of concurrent workers |
| `--timeout` | | `10s` | Per-request timeout (e.g. `5s`, `500ms`) |
| `--verbose` | `-v` | `false` | Show debug output (errors, skipped entries) |
| `--quiet` | `-q` | `false` | Suppress all output except results |
| `--status-codes` | | *(all except 404)* | Comma-separated codes to report (e.g. `200,301,302`) |
| `--follow-redirect` | | `false` | Follow HTTP redirects |
| `--user-agent` | | `goprobe/1.0` | Custom `User-Agent` header |
| `--header` | `-H` | | Custom header `"Key: Value"` (repeatable) |
| `--cookies` | | | Cookie string, e.g. `"session=abc; token=xyz"` |
| `--insecure` | `-k` | `false` | Skip TLS certificate verification |

---

## Mode Reference

### `dir` – Directory & File Brute-Force

```
goprobe dir -u <URL> -w <wordlist> [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--url` | `-u` | *(required)* | Target URL |
| `--wordlist` | `-w` | *(required)* | Path to wordlist file |
| `--extensions` | `-x` | | File extensions to append, comma-separated (e.g. `php,html,txt`) |
| `--add-slash` | | `false` | Also probe each word with a trailing `/` |
| `--method` | `-m` | `GET` | HTTP method |
| `--data` | | | POST body |

**Examples:**
```bash
# Basic directory scan
goprobe dir -u https://example.com -w /usr/share/wordlists/dirb/common.txt

# Scan for PHP and HTML files, 50 threads
goprobe dir -u https://example.com -w common.txt -x php,html -t 50

# Filter to specific response codes, with a custom header
goprobe dir -u https://example.com -w common.txt \
  --status-codes 200,301,302 \
  -H "X-Forwarded-For: 127.0.0.1"

# POST-method scan with auth cookie
goprobe dir -u https://example.com -w dirs.txt \
  -m POST --data "action=check" \
  --cookies "session=deadbeef"
```

---

### `dns` – Subdomain Enumeration

```
goprobe dns -d <domain> -w <wordlist> [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--domain` | `-d` | *(required)* | Target domain |
| `--wordlist` | `-w` | *(required)* | Path to wordlist file |
| `--resolver` | | *(system)* | Custom DNS resolver `IP[:port]` |

**Examples:**
```bash
# Basic subdomain scan
goprobe dns -d example.com -w subdomains-top1million.txt

# 100 threads, custom resolver (bypasses local cache)
goprobe dns -d example.com -w subdomains.txt -t 100 --resolver 8.8.8.8

# Use Cloudflare DNS
goprobe dns -d example.com -w subdomains.txt --resolver 1.1.1.1:53
```

---

### `vhost` – Virtual Host Discovery

```
goprobe vhost -u <URL> -w <wordlist> [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--url` | `-u` | *(required)* | Target URL/IP |
| `--wordlist` | `-w` | *(required)* | Path to wordlist file |
| `--append-domain` | | | Domain suffix to append to each word |

A baseline response size is measured before scanning. Responses that differ
from the baseline are flagged as potential vhosts.

**Examples:**
```bash
# Basic vhost scan against an IP
goprobe vhost -u http://10.10.10.1 -w vhosts.txt

# Auto-append a domain suffix
goprobe vhost -u http://10.10.10.1 -w words.txt --append-domain internal.corp.com

# Only report 200 responses
goprobe vhost -u http://10.10.10.1 -w words.txt --status-codes 200
```

---

### `s3` – AWS S3 Bucket Discovery

```
goprobe s3 -w <wordlist> [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--wordlist` | `-w` | *(required)* | Path to wordlist file |
| `--region` | | `s3.amazonaws.com` | AWS region or custom S3 endpoint |
| `--scheme` | | `https` | `https` or `http` |

**Examples:**
```bash
goprobe s3 -w buckets.txt
goprobe s3 -w buckets.txt --region us-west-2
goprobe s3 -w buckets.txt --region eu-central-1 -t 50
```

---

### `gcs` – Google Cloud Storage Bucket Discovery

```
goprobe gcs -w <wordlist> [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--wordlist` | `-w` | *(required)* | Path to wordlist file |
| `--scheme` | | `https` | `https` or `http` |

**Examples:**
```bash
goprobe gcs -w buckets.txt
goprobe gcs -w buckets.txt -t 30 --scheme http
```

---

### `tftp` – TFTP File Enumeration

```
goprobe tftp -H <host> -w <wordlist> [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--host` | `-H` | *(required)* | TFTP server hostname or IP |
| `--wordlist` | `-w` | *(required)* | Path to wordlist file |
| `--port` | | `69` | TFTP UDP port |

**Examples:**
```bash
goprobe tftp -H 192.168.1.1 -w tftp-files.txt
goprobe tftp -H 192.168.1.1 -w tftp-files.txt --port 6969 -t 20
```

---

### `fuzz` – Generic HTTP Fuzzer

```
goprobe fuzz -u <URL-with-FUZZ> -w <wordlist> [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--url` | `-u` | *(required)* | URL containing `FUZZ` token |
| `--wordlist` | `-w` | *(required)* | Path to wordlist file |
| `--method` | `-m` | `GET` | HTTP method |
| `--data` | | | POST/PUT body (may contain `FUZZ`) |
| `--fuzz-keyword` | | `FUZZ` | Custom placeholder token |

The `FUZZ` token can appear in the URL, headers (`-H`), cookies (`--cookies`),
or POST body (`--data`). All occurrences are replaced per request.

**Examples:**
```bash
# Path fuzzing
goprobe fuzz -u https://example.com/FUZZ -w wordlist.txt

# Query string parameter fuzzing
goprobe fuzz -u "https://example.com/search?q=FUZZ" -w payloads.txt

# POST body fuzzing (login brute-force)
goprobe fuzz -u https://example.com/login -m POST \
  --data "user=admin&pass=FUZZ" -w passwords.txt

# Header fuzzing (JWT token brute-force)
goprobe fuzz -u https://example.com/api/data -w tokens.txt \
  -H "Authorization: Bearer FUZZ"

# Custom keyword
goprobe fuzz -u https://example.com/INJECT -w sqli.txt \
  --fuzz-keyword INJECT --status-codes 200,500
```

---

## Output

Results are prefixed with `[+]` and written to stdout. Progress is shown inline
on TTY terminals and suppressed when piped. Pipe to `tee` to capture output:

```bash
goprobe dir -u https://example.com -w common.txt | tee results.txt
```

Use `--quiet` to suppress everything except `[+]` found lines:

```bash
goprobe dir -u https://example.com -w common.txt -q | grep '200'
```

---

## Performance Tips

- Start with `--threads 10` (default) and scale up for fast targets.
- Use `--timeout 5s` for local networks; `--timeout 15s` for slow targets.
- Use `--status-codes` to reduce noise on targets with custom error pages.
- The `dns` mode supports 100+ threads safely since it is purely network-bound.
- For very large wordlists the tool streams from disk, so RAM usage stays flat.
