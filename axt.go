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
	"time"

	"github.com/pterm/pterm"
	"github.com/tidwall/pretty"
)

// formatLevel formts the log level
//
// Returns:
// - uppercased and colorized string
// - pterm color for further use
func formatLevel(level string, emoji bool) (string, pterm.Color) {
	upperLevel := strings.ToUpper(level)
	color := pterm.FgWhite

	if emoji {
		switch upperLevel {
		case "TRACE":
			color = pterm.FgBlue
			upperLevel = "🐾 "
		case "DEBUG":
			color = pterm.FgGreen
			upperLevel = "🦠 "
		case "INFO":
			color = pterm.FgDefault
			upperLevel = "ℹ️ "
		case "WARNING", "WARN":
			color = pterm.FgYellow
			upperLevel = "⚠️ "
		case "ERROR", "ERR":
			color = pterm.FgRed
			upperLevel = "❌ "
		case "CRITICAL", "FATAL":
			color = pterm.FgMagenta
		default:
			color = pterm.FgWhite
		}
	} else {
		switch upperLevel {
		case "TRACE":
			color = pterm.FgBlue
			upperLevel = "TRACE"
		case "DEBUG":
			color = pterm.FgGreen
			upperLevel = "DEBUG"
		case "INFO":
			color = pterm.FgDefault
			upperLevel = "INFO "
		case "WARNING", "WARN":
			color = pterm.FgYellow
			upperLevel = "WARN "
		case "ERROR", "ERR":
			color = pterm.FgRed
			upperLevel = "ERR  "
		case "CRITICAL", "FATAL":
			color = pterm.FgMagenta
		default:
			color = pterm.FgWhite
		}
	}

	formattedLevel := pterm.Color(color).Sprint(upperLevel)

	return formattedLevel, color
}

// formatTime reformats time according to incoming `format` and output
func formatTime(timeStr string, format string, output string) string {
	formats := map[string]string{
		"RFC3339":     time.RFC3339,
		"RFC3339Nano": time.RFC3339Nano,
		"ANSIC":       time.ANSIC,
		"UnixDate":    time.UnixDate,
		// "2006-01-02T15:04:05.999Z07:00",
		// "2006-01-02 15:04:05.999",
		// time.RubyDate,
		// time.RFC822,
		// time.RFC822Z,
		// time.RFC850,
		// time.RFC1123,
		// time.RFC1123Z,
		// time.Kitchen,
		// time.Stamp,
		// time.StampMilli,
		// time.StampMicro,
		// time.StampNano,
	}

	chosenFormat, ok := formats[format]
	if !ok {
		return timeStr
	}

	t, err := time.Parse(chosenFormat, timeStr)
	if err != nil {
		return timeStr
	}

	return fmt.Sprintf(output, t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000)
}

// formatValue formats the value based on its type
func formatValue(value any) string {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return pterm.FgRed.Sprint(fmt.Sprintf("%v", value))
	}
	beautiful := string(pretty.Color(pretty.Pretty(jsonBytes), jsonColor()))
	return strings.TrimSuffix(beautiful, "\n")
}

// jsonColor returns a pretty style for JSON output
func jsonColor() *pretty.Style {
	return &pretty.Style{
		Key:      [2]string{"\x1B[1m\x1B[90m", "\x1B[0m"},
		String:   [2]string{"\x1B[32m", "\x1B[0m"},
		Number:   [2]string{"\x1B[33m", "\x1B[0m"},
		True:     [2]string{"\x1B[36m", "\x1B[0m"},
		False:    [2]string{"\x1B[36m", "\x1B[0m"},
		Null:     [2]string{"\x1B[2m", "\x1B[0m"},
		Escape:   [2]string{"\x1B[35m", "\x1B[0m"},
		Brackets: [2]string{"\x1B[1m", "\x1B[0m"},
		Append: func(dst []byte, c byte) []byte {
			if c < ' ' && (c != '\r' && c != '\n' && c != '\t' && c != '\v') {
				dst = append(dst, "\\u00"...)
				dst = append(dst, hexp((c>>4)&0xF))
				return append(dst, hexp((c)&0xF))
			}
			return append(dst, c)
		},
	}
}

// hexp converts a byte to its hexadecimal representation
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
