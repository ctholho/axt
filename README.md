![header-light-only](https://github.com/user-attachments/assets/800c7b87-55d1-40f7-9cb4-2ea6f1a32be5#gh-light-mode-only)
![header-dark-only](https://github.com/user-attachments/assets/9d121acb-de3f-4bfb-9a1c-b9bf6ce8063a#gh-dark-mode-only)
# Axt

axt turns structured logs to something readable on your local dev machine.

[![Go Reference](https://pkg.go.dev/badge/github.com/ctholho/axt.svg)](https://pkg.go.dev/github.com/ctholho/axt)
[![Go Report Card](https://goreportcard.com/badge/github.com/ctholho/axt)](https://goreportcard.com/report/github.com/ctholho/axt)

<hr>

Structured logs are great. But ugh, what a mess:

```
{"time":"2025-08-24T21:51:45.549605+02:00","level":"INFO","msg":"API request completed","status_code":200,"response_time_ms":127,"response_size_kb":24.7}
{"time":"2025-08-24T21:51:45.549671+02:00","level":"INFO","msg":"feature-flag:set","dark_mode":true,"beta_features":false,"thingy":null}
```

Your server app shouldn't concern itself with making logs useable on dev
machines. But you really want something like this:

```
 21:53:28.717 INFO  API request completed
                  status_code: 200
                  response_time_ms: 127
                  response_size_kb: 24.7

 21:53:28.717 INFO  feature-flag:set
                  dark_mode: true
                  beta_features: false
                  thingy: null
```

(preferably with some colors sprinkled in)

And you want to highlight if someone leaves a nasty unstructured log in:

```
🪵  something without proper JSON
```

Well, you gotta use axt (it's German for ax, btw. ...and a cool way to spell
`axed`)

## How to use

Start your application and pipe its output through axt. E.g.:

```bash
./my-application | axt
```

Your logging lib might use a variety of property names for the three important
parts of a log (as far as axt is concerned)

The defaults (as specified by go's slog lib) are:

  - level
  - msg
  - time

But you can change them with these run-time flags.

```
  -l, --level string       define name of the level property (default "level")
  -m, --message string     define name of the message property (default "msg")
  -t, --time string        define name of the time property (default "time")
```

Additionally, you can configure your output with these options:
```
```
```
  --time-in string     given time format used by time property. Uses go's time convention (default "RFC3339")
  --time-out string    print time in this format. Uses go's time convention (default "15:04:05.000")
  --emoji              display levels as emoji instead of text
  --linebreak string   "always" | only after "json" | "never" ()
```

For defining your own time input and output format refer to the go documentation of the [time format module](https://go.dev/src/time/format.go)

Examples:

- For an app using Open Telemetry:
  ```bash
  ./otel-app | axt -m EventName -t Timestamp -l SeverityText
  ```

- For Elastic Common Schema (ECS):
  ```bash
  ./spring-app | axt -t @timestamp -l log.level
  ```

- Only hours, minutes and micro seconds
  ```bash
  air | axt -time-out 04:05:000
  ```

Protips:

- Alias axt with your runtime flags to a command that makes it shorter to use
  ```bash
  "alias and-my-axt=axt -m EventName -t Timestamp -l SeverityText --time-in Unix --time-out 15:04:05.000000 --emoji --linebreak never"
  ```
  ```
  ```

- Add the pipe through axt in a Makefile of your project for the benefit of your colleagues

## Install

### Go install (recommended)

```bash
go install github.com/ctholho/axt@latest
```

If GOPATH is not part of your PATH, add axt globally for users of zsh:

```
echo "alias axt=$(go env GOPATH)/bin/axt" >> ~/.zshrc
```

### Binaries

Builds in the "Releases" section of this repository are built with the SLSA
Level 3 Compliant Go builder, which comes with strong guarantees that nobody
tampered with this build, **if you verify the build.**

Head over to the Releases section and download the binary for your system.

```bash
chmod +x ./axt-darwin-arm64
# move it to wherever you have your binaries
sudo mv axt-darwin-arm64 /usr/local/bin
```

After trying to run axt on macos for the first time, you need to head over to the
Security settings and approve axt.

## Verify Provenance

[ For up to date information head over to [SLSA verifier](https://github.com/slsa-framework/slsa-verifier#available-options) ]

To verify builds follow these steps:

First, download the matching `intoto.jsonl` file for your downloaded binary.

Next, install the slsa-verifier:

```bash
$ go install github.com/slsa-framework/slsa-verifier/v2/cli/slsa-verifier@v2.7.1
```

Finally, use slsa-verifier for axt like this:

```bash
slsa-verifier verify-artifact <PATH TO BINARY> \
  --provenance-path <PATH TO INTOTO.JSONL> \
  --source-uri github.com/ctholho/axt \
  --source-tag <VERSION NUMBER WITH LEADING v>

# For example:
slsa-verifier verify-artifact axt-darwin-arm64 \
  --provenance-path axt-darwin-arm64.intoto.jsonl \
  --source-uri github.com/ctholho/axt \
  --source-tag v0.5.2
```

## Roadmap

- [x] Allow formatting time output more freely
- [ ] Hide properties with a `--hide` option
- [ ] Color theming
- [ ] Property highlighting
