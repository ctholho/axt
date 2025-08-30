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

// levelInfo holds the display properties for a specific log level.
type levelInfo struct {
	Color pterm.Color
	Emoji string
	Text  string
}

// levelMap maps uppercase log level strings to their display properties.
// It also includes common aliases like "WARN" for "WARNING".
var levelMap = map[string]levelInfo{
	"TRACE":    {Color: pterm.FgBlue, Emoji: "🐾 ", Text: "TRACE"},
	"DEBUG":    {Color: pterm.FgGreen, Emoji: "🦠 ", Text: "DEBUG"},
	"INFO":     {Color: pterm.FgDefault, Emoji: "ℹ️ ", Text: "INFO "},
	"WARNING":  {Color: pterm.FgYellow, Emoji: "⚠️ ", Text: "WARN "},
	"WARN":     {Color: pterm.FgYellow, Emoji: "⚠️ ", Text: "WARN "},
	"ERROR":    {Color: pterm.FgRed, Emoji: "❌ ", Text: "ERR  "},
	"ERR":      {Color: pterm.FgRed, Emoji: "❌ ", Text: "ERR  "},
	"CRITICAL": {Color: pterm.FgMagenta, Emoji: "❌ ", Text: "CRITICAL"},
	"FATAL":    {Color: pterm.FgMagenta, Emoji: "❌ ", Text: "FATAL"},
}

// formatLevel formts the log level
//
// Returns:
// - uppercased and colorized string
// - pterm color for further use.
func formatLevel(level string, useEmoji bool) (string, pterm.Color) {
	upperLevel := strings.ToUpper(level)

	if info, ok := levelMap[upperLevel]; ok {
		formattedLevel := info.Color.Sprint(info.Text)
		if useEmoji && info.Emoji != "" {
			formattedLevel = info.Emoji
		}

		return formattedLevel, info.Color
	}

	return upperLevel, pterm.FgWhite
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
		Key:      [2]string{"\x1B[1m\x1B[90m", "\x1B[0m"},
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
