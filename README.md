![header-light-only](https://github.com/user-attachments/assets/800c7b87-55d1-40f7-9cb4-2ea6f1a32be5#gh-light-mode-only)
![header-dark-only](https://github.com/user-attachments/assets/9d121acb-de3f-4bfb-9a1c-b9bf6ce8063a#gh-dark-mode-only)
# axt

Turn structured logs to something readable on your local dev machine.

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

<table align="center">
  <tr>
    <td align="center">
      <img width="469" height="483" alt="axt-hide-time" src="https://github.com/user-attachments/assets/2e3915b1-24ae-449c-9dbc-d181fd8b2fde" />
      <br>
      <em>Hide timestamps with `--hide time`</em>
    </td>
    <td align="center">
      <img width="469" height="483" alt="axt-emoji" src="https://github.com/user-attachments/assets/a406e7ae-cc74-40fb-a8a6-3354dfd3e631" />
      <br>
      <em>Use emojis for log levels with `--emoji`</em>
    </td>
  </tr>
</table>

And you want to highlight if someone leaves in a nasty unstructured log:

```
🪵  something without proper JSON
```

Well, you gotta use axt (it's German for axe, btw. ...and a cool way to spell
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
  --time-in string     given time format used by time property. Uses go's time convention; or use 'Unix' | 'UnixMilli' | 'UnixMicro' for Unix epoch timestamps.(default "RFC3339")
  --time-out string    print time in this format. Uses go's time convention (default "15:04:05.000")
  --emoji              display levels as emoji instead of text
  --linebreak string   "always" | only after "json" | "never" (default: always)
  --hide string        hide a property. Use the flag multiple times to hide more than one.
```

For defining your own time input and output format refer to the go documentation of the [time format module](https://go.dev/src/time/format.go)

Examples:

- Using Open Telemetry:
  `./otel-app | axt -m EventName -t Timestamp -l SeverityText`

- Elastic Common Schema:
  ` ./spring-app | axt -t @timestamp -l log.level`

- Only hours, minutes and micro seconds
  `air | axt -time-out 04:05:000`

- Hide time
  `timeless-app | axt --hide time`

Protips:

- Alias axt with your runtime flags to a command that makes it shorter to use
  ```bash
  "alias and-my-axt=axt -m EventName -t Timestamp -l SeverityText --time-in Unix --time-out 15:04:05.000000 --emoji --linebreak never"
  ```

- Add the pipe through axt in a Makefile of your project for the benefit of your colleagues

## Install

### Brew

Install with

```bash
brew tap ctholho/ctholho
brew install axt
```

See below, in [Verify Provenance](#verify-provenance) how to verify your build
with SLSA Level 3 compliance guarantees.

### Go install

```bash
go install github.com/ctholho/axt@latest
```

If GOPATH is not part of your PATH, add axt globally for users of zsh:

```
echo "alias axt=$(go env GOPATH)/bin/axt" >> ~/.zshrc
```

## Verify Provenance

Builds in the "Releases" section of this repository are built with the SLSA
Level 3 Compliant Go builder, which comes with strong guarantees that nobody
tampered with this build, **if you verify the build.**

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
slsa-verifier verify-artifact /opt/homebrew/bin/axt \
  --provenance-path ~/Downloads/axt-darwin-arm64.intoto.jsonl \
  --source-uri github.com/ctholho/axt \
  --source-tag v0.7.0
```

## Development

As a handy helper during development use `debug/debug.go` to simulate a server.

1. Build the axt binary with `go build -o axt`
2. Use the binary on the debug file: `go run debug/debug.go | ./axt

## Roadmap

- [ ] Color theming
- [ ] Property highlighting
