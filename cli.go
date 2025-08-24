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

var (
	msgKeyDefault            = "msg"
	timeKeyDefault           = "time"
	levelKeyDefault          = "level"
	emptyLineStrategyDefault = "always"
	emojiLevelDefault        = false
	timeInputFormatDefault   = "RFC3339"
	timeOutputFormatDefault  = "%02d:%02d:%02d.%03d"
)

func main() {
	timeKeyFlag := flag.StringP("time", "t", timeKeyDefault, "define name of the time property")
	messageKeyFlag := flag.StringP("message", "m", msgKeyDefault, "define name of the message property")
	levelKeyFlag := flag.StringP("level", "l", levelKeyDefault, "define name of the level property")
	emptyLineStrategyFlag := flag.String("linebreak", emptyLineStrategyDefault, "\"always\" | only after \"json\" | \"never\"")
	emojiLevel := flag.Bool("emoji", emojiLevelDefault, "display levels as emoji instead of text")
	timeInputFormatFlag := flag.String("time-in", timeInputFormatDefault, "given time format used by time property. Some values used by go's time module are possible.")
	timeOutputFormatFlag := flag.String("time-out", timeOutputFormatDefault, "print time in this format. (WIP!)")

	var showVersion bool
	flag.BoolVarP(&showVersion, "version", "v", false, "Show version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "axt | structured logs but forcibly gemütlich | %s\n\n", Version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  axt [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			fmt.Printf("🪵  %s\n%s", line, formatNewLine(*emptyLineStrategyFlag, false))
			continue
		}

		// TIME
		timeValue, _ := entry[*timeKeyFlag].(string)
		timeColor := pterm.FgDarkGray
		t := formatTime(timeValue, *timeInputFormatFlag, *timeOutputFormatFlag)
		formattedTime := pterm.Color(timeColor).Sprint(t)

		// LEVEL
		levelValue, ok := entry[*levelKeyFlag].(string)
		if levelValue == "" || !ok {
			// Should never happen, but who knows how people configure their logger
			levelValue = "NO LEVEL"
		}
		formattedLevel, levelColor := formatLevel(levelValue, *emojiLevel)

		// MESSAGE
		messageValue, _ := entry[*messageKeyFlag].(string)
		formattedMessage := pterm.Color(levelColor).Sprint(messageValue)

		// OVERALL FORMAT
		fmt.Printf(" %s %s %s\n", formattedTime, formattedLevel, formattedMessage)

		// Remove standard fields
		delete(entry, *timeKeyFlag)
		delete(entry, *messageKeyFlag)
		delete(entry, *levelKeyFlag)

		lineColor := pterm.FgGray
		verticalLine := pterm.Color(lineColor).Sprint("               ") // Add alignment for short events

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

		// Show a pretty vertical line if the event is longer
		if len(logLines) > 3 {
			outerColor := pterm.NewRGB(70, 70, 70)
			innerColor := pterm.NewRGB(150, 150, 150)

			// Fade for the line
			for i, line := range logLines {
				var currentVerticalLine string
				if i == 0 {
					currentVerticalLine = outerColor.Fade(0, float32(len(logLines)-1), float32(i), innerColor, outerColor).Sprint("              ┌")
				} else if i == len(logLines)-1 {
					currentVerticalLine = outerColor.Fade(0, float32(len(logLines)-1), float32(i), innerColor, outerColor).Sprint("              └")
				} else {
					currentVerticalLine = outerColor.Fade(0, float32(len(logLines)-1), float32(i), innerColor, outerColor).Sprint("              │")
				}
				fmt.Println(strings.Replace(line, verticalLine, currentVerticalLine, 1))
			}
		} else {
			for _, line := range logLines {
				fmt.Println(line)
			}
		}

		// Maybe an empty line after each event
		fmt.Printf("%s", formatNewLine(*emptyLineStrategyFlag, true))

	}
	if err := scanner.Err(); err != nil {
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
