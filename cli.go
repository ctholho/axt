package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/pterm/pterm"
	flag "github.com/spf13/pflag"
)

var (
	Version    = "dev"
	Commit     = "local"
	CommitDate = "n/a"
	TreeState  = "dirty"
)

func main() {
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

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		var entry map[string]any

		err := json.Unmarshal([]byte(line), &entry)
		if err != nil {
			fmt.Printf("🪵  %s\n%s", line, formatNewLine(cfg.EmptyLineStrategy, false))

			continue
		}

		// TIME
		timeValue, _ := entry[cfg.TimeKey].(string)
		timeColor := pterm.FgDarkGray
		t := formatTime(timeValue, cfg.TimeInputFormat, cfg.TimeOutputFormat)
		formattedTime := timeColor.Sprint(t)

		// LEVEL
		levelValue, ok := entry[cfg.LevelKey].(string)
		if levelValue == "" || !ok {
			levelValue = "NO LEVEL"
		}

		formattedLevel, levelColor := formatLevel(levelValue, cfg.EmojiLevel)

		// MESSAGE
		messageValue, _ := entry[cfg.MessageKey].(string)
		formattedMessage := levelColor.Sprint(messageValue)

		// OVERALL FORMAT of first line
		fmt.Printf(" %s %s %s\n", formattedTime, formattedLevel, formattedMessage)

		// Remove standard fields to avoid duplication
		delete(entry, cfg.TimeKey)
		delete(entry, cfg.MessageKey)
		delete(entry, cfg.LevelKey)

		lineColor := pterm.FgGray
		verticalLine := lineColor.Sprint("               ") // Add alignment for short events

		var logLines []string

		// add extra fields if any
		if len(entry) > 0 {
			for key, value := range entry {
				formattedKey := pterm.NewStyle(pterm.FgWhite, pterm.Bold).Sprint(key)
				formattedValue := formatValue(value)
				formattedValueLines := strings.Split(formattedValue, "\n")
				logLines = append(logLines, fmt.Sprintf("%s   %s: %s", verticalLine, formattedKey, formattedValueLines[0]))

				for _, line := range formattedValueLines[1:] {
					logLines = append(logLines, fmt.Sprintf("%s   %s", verticalLine, line))
				}
			}
		}

		// Show a pretty vertical line if there's some properties (at least 3)
		addBorder(logLines, verticalLine)

		// Maybe add an empty line after each event
		fmt.Printf("%s", formatNewLine(cfg.EmptyLineStrategy, true))
	}

	err := scanner.Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
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

type Config struct {
	TimeKey           string
	MessageKey        string
	LevelKey          string
	EmptyLineStrategy string
	EmojiLevel        bool
	TimeInputFormat   string
	TimeOutputFormat  string
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
	}
}

func setupFlags(cfg *Config) {
	flag.StringVarP(&cfg.TimeKey, "time", "t", cfg.TimeKey, "define name of the time property")
	flag.StringVarP(&cfg.MessageKey, "message", "m", cfg.MessageKey, "define name of the message property")
	flag.StringVarP(&cfg.LevelKey, "level", "l", cfg.LevelKey, "define name of the level property")
	flag.StringVar(&cfg.EmptyLineStrategy, "linebreak", cfg.EmptyLineStrategy, "\"always\" | only after \"json\" | \"never\"")
	flag.BoolVar(&cfg.EmojiLevel, "emoji", cfg.EmojiLevel, "display levels as emoji instead of text")
	flag.StringVar(
		&cfg.TimeInputFormat,
		"time-in",
		cfg.TimeInputFormat,
		`format of time property. Uses go's time convention; or use 'Unix' | 'UnixMilli' | 'UnixMicro' for Epoch timestamps.`,
	)
	flag.StringVar(&cfg.TimeOutputFormat, "time-out", cfg.TimeOutputFormat, "print time in this format. Uses go's time format convention.")
}
