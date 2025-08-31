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
	Style     *pterm.Style
	MainColor pterm.Color
	Emoji     string
	Text      string
	Align     bool
}

// levelMap maps uppercase log level strings to their display properties.
// It also includes common aliases like "WARN" for "WARNING".
var levelMap = map[string]levelInfo{
	"TRACE":    {Align: true, Style: pterm.NewStyle(pterm.BgBlue, pterm.FgBlack, pterm.Bold), MainColor: pterm.FgBlue, Emoji: "🐾 ", Text: " TRACE "},
	"DEBUG":    {Align: true, Style: pterm.NewStyle(pterm.BgGreen, pterm.FgBlack, pterm.Bold), MainColor: pterm.FgGreen, Emoji: "🦠 ", Text: " DEBUG "},
	"INFO":     {Align: false, Style: pterm.NewStyle(pterm.BgBlue, pterm.FgBlack, pterm.Bold), MainColor: pterm.FgDefault, Emoji: "ℹ️ ", Text: "  INFO  "},
	"WARN":     {Align: false, Style: pterm.NewStyle(pterm.BgYellow, pterm.FgBlack, pterm.Bold), MainColor: pterm.FgYellow, Emoji: "⚠️ ", Text: "  WARN  "},
	"WARNING":  {Align: true, Style: pterm.NewStyle(pterm.BgYellow, pterm.FgBlack, pterm.Bold), MainColor: pterm.FgYellow, Emoji: "⚠️ ", Text: "WARNING"},
	"ERROR":    {Align: true, Style: pterm.NewStyle(pterm.BgRed, pterm.FgBlack, pterm.Bold), MainColor: pterm.FgRed, Emoji: "❌ ", Text: " ERROR "},
	"ERR":      {Align: true, Style: pterm.NewStyle(pterm.BgRed, pterm.FgBlack, pterm.Bold), MainColor: pterm.FgRed, Emoji: "❌ ", Text: "  ERR  "},
	"FATAL":    {Align: true, Style: pterm.NewStyle(pterm.BgRed, pterm.FgBlack, pterm.Bold), MainColor: pterm.FgMagenta, Emoji: "❌ ", Text: " FATAL "},
	"CRITICAL": {Align: false, Style: pterm.NewStyle(pterm.BgRed, pterm.FgBlack, pterm.Bold), MainColor: pterm.FgMagenta, Emoji: "❌ ", Text: "CRITICAL"},
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
			formattedLevel = fmt.Sprintf("%s ", formattedLevel)
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
