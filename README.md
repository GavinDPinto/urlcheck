# urlcheck

A small command-line tool that checks a list of URLs in parallel and reports which ones are up.

Built as practice for [app-operator](https://github.com/GavinDPinto/app-operator) (Phase 0, issue #2).

## What it does

```
$ urlcheck -file urls.txt -workers 5 -timeout 3s
UP    200   142ms  https://go.dev
UP    200   310ms  https://github.com
DOWN  500    95ms  https://httpbin.org/status/500
DOWN  -    3001ms  https://httpbin.org/delay/10       (timeout)
DOWN  -       4ms  https://does-not-exist.invalid     (no such host)

3 of 5 down
```

- Reads URLs from a file, one per line. Blank lines and lines starting with `#` are skipped.
- Sends an HTTP GET to each URL, several at a time.
- A URL is **UP** if it answers with a status code in the 200–399 range, and **DOWN** if it gets any other status, an error, or no answer before the timeout.
- Prints one line per URL with status, status code, and response time, then a summary.
- Exits with code `1` if anything is down and `0` if everything is up, so scripts and CI can use it.

### Flags

| Flag | Default | Meaning |
| --- | --- | --- |
| `-file` | `urls.txt` | File to read URLs from |
| `-workers` | `5` | How many URLs to check at the same time |
| `-timeout` | `5s` | How long to wait for each URL |

## Milestones

Work through these in order. Each one should run before you start the next.

1. **Read the file.** Open `urls.txt`, read it line by line, skip blanks and `#` comments, and print each URL.
   _Look up:_ `os.Open`, `bufio.Scanner`, `strings.TrimSpace`, `strings.HasPrefix`.

2. **Check one URL.** Write `check(url string) Result` that does a GET and returns a `Result` struct (URL, status code, duration, error). Use an `http.Client` with a `Timeout`. Close the response body.
   _Look up:_ `net/http` (`http.Client`, `resp.StatusCode`, `resp.Body.Close`), `time.Now`, `time.Since`.

3. **Check all of them, one by one.** Loop over the URLs, call `check`, and print each result. Note how long the whole run takes.

4. **Check them all at once.** Start one goroutine per URL. Each one sends its `Result` on a channel, and `main` receives and prints them. This time use **channels** for the results instead of a shared map with a mutex. Compare the total time with step 3.
   _Look up:_ `sync.WaitGroup` again, and the "close the channel when all workers are done" pattern.

5. **Limit how many run at once.** With 1,000 URLs you wouldn't want 1,000 requests in flight at the same time. Use a fixed number of worker goroutines reading URLs from a jobs channel (a *worker pool*).

6. **Flags, summary, exit code.** Add `-file`, `-workers` and `-timeout` with the `flag` package. Print the "N of M down" summary, and exit with `os.Exit(1)` if anything is down.

### Stretch goals

- Print the results in the same order as the file, even though they finish in a random order.
- A `-json` flag that prints the results as JSON.
- Retry a failed URL once before marking it DOWN.

## Running

```
go run . -file urls.txt
```

## Tests

Issue #3: unit tests using `net/http/httptest` to start fake servers that return 200, return 500, or never answer.
