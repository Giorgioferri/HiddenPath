# minigobuster

A small directory and path brute-forcer written in Go, built for learning offensive security and Go's concurrency model. Given a target URL and a wordlist, it requests each path and reports the ones that exist, based on the HTTP status code.

This is an **educational** project — a stripped-down clone of [gobuster](https://github.com/OJ/gobuster), written from scratch to practice file I/O, the `net/http` client, and goroutines.

> ⚠️ Only run this against targets you own or are explicitly authorized to test (your own local server, lab machines, TryHackMe / HackTheBox boxes). Scanning systems you don't have permission for can be illegal.

## How it works

The tool reads a wordlist line by line. For each word it builds a URL (`site/word`) and sends an HTTP `GET` request, then looks at the **status code** of the response:

- `200` — the path exists and is accessible
- `301` / `302` — redirect, usually means a directory exists
- `403` — forbidden, but the path exists (often interesting)
- `404` — not found, skipped

Anything that is **not** `404` is reported, so protected (`403`) and redirected (`301`) paths aren't missed.

Requests are sent concurrently using goroutines, so the whole wordlist is scanned in parallel instead of one request at a time, with a `sync.WaitGroup` making sure the program waits for every request to finish before exiting.

## Requirements

- Go 1.22 or newer

## Build

```bash
git clone https://github.com/Giorgioferri/minigobuster.git
cd minigobuster
go build -o minigobuster .
```

Or run it directly without building:

```bash
go run . -site http://localhost:8000 -wordlist wordlist.txt
```

## Usage

```bash
minigobuster -site <url> -wordlist <path>
```

| Flag | Description | Default |
|------|-------------|---------|
| `-site` | Target base URL (include `http://`) | `url` |
| `-wordlist` | Path to the wordlist file | `wordlist.txt` |

Example with a real wordlist:

```bash
minigobuster -site http://localhost:8000 -wordlist /usr/share/wordlists/dirb/common.txt
```

## Example output

```
found! [301] http://localhost:8000/admin
found! [200] http://localhost:8000/secret.txt
```

## Testing it locally

You can test the tool safely against a local web server you control. Create a folder with a couple of known subfolders/files, then serve it:

```bash
mkdir testsite
cd testsite
mkdir admin
echo ok > secret.txt
python -m http.server 8000
```

Then point the tool at `http://localhost:8000` with a wordlist that contains `admin` and `secret.txt`. Because you know exactly what exists, you can confirm the output is correct — `admin` should come back as `301`, `secret.txt` as `200`, and anything else as `404`.

## Roadmap

- [ ] **Proper concurrency with a worker pool** — the current version launches one goroutine per word. Cap the number of simultaneous requests (e.g. 50 workers) using buffered channels, so large wordlists don't open thousands of connections at once and overwhelm the machine or the target.
- [ ] **`-timeout` flag** — switch from `http.Get` to a configurable `http.Client{Timeout}` so a single hanging request can't freeze the whole scan.
- [ ] **`-ext` flag** — try multiple file extensions per word (e.g. `admin.php`, `admin.html`, `admin.bak`).
- [ ] **`-output` flag** — save the discovered paths to a file.

## License

MIT — feel free to use and modify it.
