// 🌈 pretty Print Logs from afm-core 🌈
// Usage:
//
// ```bash
// ./your_app | ./axt
//
// # Or read from a log file
// tail -f /path/to/logfile | ./axt

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/tidwall/pretty"
)

var (
	msgKey   = "msg"
	timeKey  = "time"
	levelKey = "level"
)

// formatLevel formts the log level
//
// Returns:
// - uppercased and colorized string
// - pterm color for further use
func formatLevel(level string) (string, pterm.Color) {
	upperLevel := strings.ToUpper(level)
	color := pterm.FgWhite

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

	formattedLevel := pterm.Color(color).Sprint(upperLevel)

	return formattedLevel, color
}

// formatTime reformats time to display only seconds and milliseconds
func formatTime(timeStr string) string {
	// Try parsing in different formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		// "2006-01-02T15:04:05.999Z07:00",
		// "2006-01-02 15:04:05.999",
		// time.ANSIC,
		// time.UnixDate,
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

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, timeStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return timeStr
	}

	return fmt.Sprintf("%02d:%02d:%02d.%03d", t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000)
}

// formatValue formats the value based on its type
func formatValue(value any) string {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return pterm.FgRed.Sprint(fmt.Sprintf("%v", value))
	}
	beautiful := string(pretty.Color(pretty.Pretty(jsonBytes), JsonColor()))
	return strings.TrimSuffix(beautiful, "\n")
}

// JsonColor returns a pretty style for JSON output
func JsonColor() *pretty.Style {
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			fmt.Printf("🪵  %s\n\n", line)
			continue
		}

		// TIME
		timeValue, _ := entry[timeKey].(string)
		timeColor := pterm.FgDarkGray
		formattedTime := pterm.Color(timeColor).Sprint(formatTime(timeValue))

		// LEVEL
		levelValue, ok := entry[levelKey].(string)
		if levelValue == "" || !ok {
			// Should never happen, but who knows how people configure their logger
			levelValue = "NO LEVEL"
		}
		formattedLevel, levelColor := formatLevel(levelValue)

		// MESSAGE
		messageValue, _ := entry[msgKey].(string)
		formattedMessage := pterm.Color(levelColor).Sprint(messageValue)

		// OVERALL FORMAT
		fmt.Printf(" %s %s %s\n", formattedTime, formattedMessage, formattedLevel)

		// Remove standard fields
		delete(entry, timeKey)
		delete(entry, levelKey)
		delete(entry, msgKey)

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
		// Blank line after each event
		fmt.Println()

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
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
