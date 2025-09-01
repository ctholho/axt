package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	flag "github.com/spf13/pflag"
)

var (
	Version    = "dev"
	Commit     = "local"
	CommitDate = "n/a"
	TreeState  = "dirty"
)

type Config struct {
	TimeKey           string
	MessageKey        string
	LevelKey          string
	EmptyLineStrategy string
	EmojiLevel        bool
	TimeInputFormat   string
	TimeOutputFormat  string
	HiddenKeys        []string
}

func newConfig() *Config {
	return &Config{
		TimeKey:           "time",
		MessageKey:        "msg",
		LevelKey:          "level",
		EmptyLineStrategy: "always",
		EmojiLevel:        false,
		TimeInputFormat:   "RFC3339",
		TimeOutputFormat:  "15:04:05.000",
		HiddenKeys:        []string{},
	}
}

func setupFlags(cfg *Config) {
	flag.StringVarP(&cfg.TimeKey, "time", "t", cfg.TimeKey, "Name of the time property")
	flag.StringVarP(&cfg.MessageKey, "message", "m", cfg.MessageKey, "Name of the message property")
	flag.StringVarP(&cfg.LevelKey, "level", "l", cfg.LevelKey, "Name of the level property")
	flag.StringVar(&cfg.EmptyLineStrategy, "linebreak", cfg.EmptyLineStrategy, "\"always\" | only after \"json\" | \"never\"")
	flag.BoolVar(&cfg.EmojiLevel, "emoji", cfg.EmojiLevel, "Display levels as emoji instead of text")
	flag.StringVar(
		&cfg.TimeInputFormat,
		"time-in",
		cfg.TimeInputFormat,
		`Go time layout string or 'Unix' | 'UnixMilli' | 'UnixMicro'.`,
	)
	flag.StringVar(&cfg.TimeOutputFormat, "time-out", cfg.TimeOutputFormat, "Print time in this format. Use Go time format string.")
	flag.StringSliceVar(&cfg.HiddenKeys,
		"hide",
		cfg.HiddenKeys,
		"Hide a property. Use the flag multiple times to hide more than one.")
}

func setupCLI() *Config {
	cfg := newConfig()
	setupFlags(cfg)

	var showVersion bool

	flag.BoolVarP(&showVersion, "version", "v", false, "Show version information")

	flag.Usage = printHelp
	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
	}

	return cfg
}

func printVersion() {
	buildInfo, ok := debug.ReadBuildInfo()
	if Version == "dev" && ok {
		Version = buildInfo.Main.Version
	}

	fmt.Printf("axt version: %s\n", Version)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Built at: %s\n", CommitDate)
	fmt.Printf("Tree state: %s\n", TreeState)
}

func printHelp() {
	fmt.Fprintf(os.Stderr, "axt | structured logs but forcibly gemütlich | %s\n\n", Version)
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  axt [options]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func scan(cfg *Config) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		var entry map[string]any

		err := json.Unmarshal([]byte(line), &entry)
		if err != nil {
			prettyPrintBadJSON(line, cfg)

			continue
		}

		prettyPrintJSON(entry, cfg)
	}

	err := scanner.Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	cfg := setupCLI()
	scan(cfg)
}
