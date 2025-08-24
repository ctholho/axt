![header-light-only](https://github.com/user-attachments/assets/800c7b87-55d1-40f7-9cb4-2ea6f1a32be5#gh-light-mode-only)
![header-dark-only](https://github.com/user-attachments/assets/9d121acb-de3f-4bfb-9a1c-b9bf6ce8063a#gh-dark-mode-only)
# Axt

axt turns structured logs to something readable on your local dev machine.
It's mostly a wrapper for the awesome project [pterm](github.com/pterm)

Basically, something like this:

```
{"time":"2025-08-24T21:51:45.549605+02:00","level":"INFO","msg":"API request completed","status_code":200,"response_time_ms":127,"response_size_kb":24.7}
{"time":"2025-08-24T21:51:45.549671+02:00","level":"INFO","msg":"Feature flags status","dark_mode":true,"beta_features":false,"thingy":null}
```

becomes:

```
 21:53:28.717 INFO  API request completed
                  status_code: 200
                  response_time_ms: 127
                  response_size_kb: 24.7

 21:53:28.717 INFO  Feature flags status
                  dark_mode: true
                  beta_features: false
                  thingy: null
```

(with some colors)

Unstructured logs, that found it's way in your logging output (purely by,
accident of course) will be displayed like

```
🪵  something without proper JSON
```

## Install

### With go install (recommended)

```bash
go install github.com/ctholho/axt@latest

# e.g. if you use zsh
alias axt="$(go env GOPATH)/bin/axt" >> ~/.zshrc

# or if you use bash
alias axt="$(go env GOPATH)/bin/axt" >> ~/.bashrc
```


### SLSA Level 3 compliant binaries

Builds in the "Releases" section of this repository are built with the SLSA
Level 3 Compliant Go builder, which comes with strong guarantees that nobody
tampered with this build if you verify the build.

Head over to the Releases section and download the binary for your system.

```bash
chmod +x ./axt-darwin-arm64
# move it to wherever you want to have your executable binaries
sudo mv axt-darwin-arm64 /usr/local/bin
```

After trying to run it on macos for the first time, you need to head over to the
Security settings and approve axt.


### Verify Provenance

To verify your build binary from the "Releases" follow these steps.

For up to date information and latest version numbers head over to
[SLSA verifier](https://github.com/slsa-framework/slsa-verifier#available-options)

Download the proper `intoto.jsonl` file for your binary.

Install the slsa-verifier

```bash
$ go install github.com/slsa-framework/slsa-verifier/v2/cli/slsa-verifier@v2.7.1
```

Then, use the slsa-verifier like this:

```bash
slsa-verifier verify-artifact <path to binary> \
  --provenance-path <path to intoto.jsonl> \
  --source-uri github.com/ctholho/axt \
  --source-tag <version number with leading v>

# For example:
slsa-verifier verify-artifact axt-darwin-arm64 \
  --provenance-path axt-darwin-arm64.intoto.jsonl \
  --source-uri github.com/ctholho/axt \
  --source-tag v0.5.2

```

## Options

In case your application uses different property names for time, msg or level
you can set these with run-time flags.

These run-time flags are available:

```
  -l, --level string       define name of the level property (default "level")
  -m, --message string     define name of the message property (default "msg")
  -t, --time string        define name of the time property (default "time")
      --time-in string     given time format used by time property. Some values used by go's time module are possible. (default "RFC3339")
      --emoji              display levels as emoji instead of text
      --linebreak string   "always" | only after "json" | "never" (default "always")
```

