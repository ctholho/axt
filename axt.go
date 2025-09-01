// Usage:
//
// ```bash
// ./your_app | ./axt
//
// # Or read from a log file
// tail -f /path/to/logfile | ./axt

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/tidwall/pretty"
)

func prettyPrintJSON(entry map[string]any, cfg *Config) {
	// Remove properties if a user wants to hide them
	hideProperties(entry, cfg.HiddenKeys...)

	// TIME
	timeValue, _ := entry[cfg.TimeKey].(string)
	timeColor := pterm.FgDarkGray
	formattedTime := formatTime(timeValue, cfg.TimeInputFormat, cfg.TimeOutputFormat)

	var formattedTimeWithAlign string
	if timeValue == "" {
		formattedTimeWithAlign = ""
	} else {
		formattedTimeWithAlign = timeColor.Sprintf(`%s `, formattedTime)
	}

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
	fmt.Printf("%s%s %s\n", formattedTimeWithAlign, formattedLevel, formattedMessage)

	// Remove standard properties to avoid duplication if we display them on the
	// first line
	keysToHide := []string{cfg.TimeKey, cfg.LevelKey, cfg.MessageKey}
	hideProperties(entry, keysToHide...)

	lineColor := pterm.FgGray
	// Add alignment
	vertAlign := lineColor.Sprint("      ")

	var logLines []string

	// add extra fields if any
	if len(entry) > 0 {
		for key, value := range entry {
			formattedKey := pterm.NewStyle(pterm.FgDefault).Sprint(key)
			formattedValue := formatValue(value)
			formattedValueLines := strings.Split(formattedValue, "\n")
			logLines = append(logLines, fmt.Sprintf("%s   %s: %s", vertAlign, formattedKey, formattedValueLines[0]))

			for _, line := range formattedValueLines[1:] {
				logLines = append(logLines, fmt.Sprintf("%s   %s", vertAlign, line))
			}
		}
	}

	// Show a pretty vertical line if there's some properties (at least 3)
	addBorder(logLines, vertAlign)

	// Maybe add an empty line after each event
	fmt.Printf("%s", formatNewLine(cfg.EmptyLineStrategy, true))
}

// levelInfo holds the display properties for a specific log level.
type levelInfo struct {
	Style     *pterm.Style
	MainColor pterm.Color
	Emoji     string
	Text      string
	Align     bool
}

// levelMap maps uppercase log level strings to their display properties.
// It also includes common aliases like "WARN" for "WARNING".
var levelMap = map[string]levelInfo{
	"TRACE": {
		Align:     true,
		Style:     pterm.NewStyle(pterm.BgBlue, pterm.FgBlack, pterm.Bold),
		MainColor: pterm.FgBlue,
		Emoji:     "🐾 ",
		Text:      " TRACE ",
	},
	"DEBUG": {
		Align:     true,
		Style:     pterm.NewStyle(pterm.BgGreen, pterm.FgBlack, pterm.Bold),
		MainColor: pterm.FgGreen,
		Emoji:     "🦠 ",
		Text:      " DEBUG ",
	},
	"INFO": {
		Align:     false,
		Style:     pterm.NewStyle(pterm.BgBlue, pterm.FgBlack, pterm.Bold),
		MainColor: pterm.FgDefault,
		Emoji:     "ℹ️ ",
		Text:      "  INFO  ",
	},
	"WARN": {
		Align:     false,
		Style:     pterm.NewStyle(pterm.BgYellow, pterm.FgBlack, pterm.Bold),
		MainColor: pterm.FgYellow,
		Emoji:     "⚠️ ",
		Text:      "  WARN  ",
	},
	"WARNING": {
		Align:     true,
		Style:     pterm.NewStyle(pterm.BgYellow, pterm.FgBlack, pterm.Bold),
		MainColor: pterm.FgYellow,
		Emoji:     "⚠️ ",
		Text:      "WARNING",
	},
	"ERROR": {
		Align:     true,
		Style:     pterm.NewStyle(pterm.BgRed, pterm.FgBlack, pterm.Bold),
		MainColor: pterm.FgRed,
		Emoji:     "❌ ",
		Text:      " ERROR ",
	},
	"ERR": {
		Align:     true,
		Style:     pterm.NewStyle(pterm.BgRed, pterm.FgBlack, pterm.Bold),
		MainColor: pterm.FgRed,
		Emoji:     "❌ ",
		Text:      "  ERR  ",
	},
	"FATAL": {
		Align:     true,
		Style:     pterm.NewStyle(pterm.BgRed, pterm.FgBlack, pterm.Bold),
		MainColor: pterm.FgMagenta,
		Emoji:     "❌ ",
		Text:      " FATAL ",
	},
	"CRITICAL": {
		Align:     false,
		Style:     pterm.NewStyle(pterm.BgRed, pterm.FgBlack, pterm.Bold),
		MainColor: pterm.FgMagenta,
		Emoji:     "❌ ",
		Text:      "CRITICAL",
	},
}

func prettyPrintBadJSON(line string, cfg *Config) {
	fmt.Printf("🪵  %s\n%s", line, formatNewLine(cfg.EmptyLineStrategy, false))
}

// formatLevel formts the log level
//
// Returns:
// - uppercased and colorized string
// - pterm color of level for further use.
func formatLevel(level string, useEmoji bool) (string, pterm.Color) {
	levelUppercase := strings.ToUpper(level)

	if info, ok := levelMap[levelUppercase]; ok {
		formattedLevel := info.Style.Sprint(info.Text)
		if useEmoji && info.Emoji != "" {
			formattedLevel = info.Emoji
		}

		if info.Align {
			formattedLevel += " "
		}

		return formattedLevel, info.MainColor
	}

	return levelUppercase, pterm.FgDefault
}

// formatValue formats the value based on its type.
func formatValue(value any) string {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return pterm.FgRed.Sprint(fmt.Sprintf("%v", value))
	}

	beautiful := string(pretty.Color(pretty.Pretty(jsonBytes), jsonColor()))

	return strings.TrimSuffix(beautiful, "\n")
}

// jsonColor returns a pretty style for JSON output.
func jsonColor() *pretty.Style {
	return &pretty.Style{
		Key:      [2]string{"\x1B[39m", "\x1B[0m"},
		String:   [2]string{"\x1B[32m", "\x1B[0m"},
		Number:   [2]string{"\x1B[33m", "\x1B[0m"},
		True:     [2]string{"\x1B[36m", "\x1B[0m"},
		False:    [2]string{"\x1B[36m", "\x1B[0m"},
		Null:     [2]string{"\x1B[2m", "\x1B[0m"},
		Escape:   [2]string{"\x1B[35m", "\x1B[0m"},
		Brackets: [2]string{"\x1B[1m", "\x1B[0m"},
		Append: func(dst []byte, cur byte) []byte {
			if cur < ' ' && (cur != '\r' && cur != '\n' && cur != '\t' && cur != '\v') {
				dst = append(dst, "\\u00"...)
				dst = append(dst, hexp((cur>>4)&0xF))

				return append(dst, hexp((cur)&0xF))
			}

			return append(dst, cur)
		},
	}
}

// hexp converts a byte to its hexadecimal representation.
func hexp(p byte) byte {
	switch {
	case p < 10:
		return p + '0'
	default:
		return (p - 10) + 'a'
	}
}

// formatNewLine returns a string that is contains either a new line or is empty.
func formatNewLine(strategy string, structured bool) string {
	switch strategy {
	case "always":
		return "\n"
	case "json":
		if structured {
			return "\n"
		}

		return ""
	case "never":
		return ""
	default:
		return ""
	}
}

// hideProperty removes keys from a map.
// The function modifies the map in place.
func hideProperties(entry map[string]any, keys ...string) {
	for _, k := range keys {
		delete(entry, k)
	}
}
