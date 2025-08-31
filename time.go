package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	unixStrategy string            = "unixStrategy"
	formats      map[string]string = map[string]string{
		"RFC3339":                       time.RFC3339,
		"RFC3339Nano":                   time.RFC3339Nano,
		"ANSIC":                         time.ANSIC,
		"UnixDate":                      time.UnixDate,
		"2006-01-02T15:04:05.999Z07:00": "2006-01-02T15:04:05.999Z07:00",
		"2006-01-02 15:04:05.999":       "2006-01-02 15:04:05.999",
		"RubyDate":                      time.RubyDate,
		"RFC822":                        time.RFC822,
		"RFC822Z":                       time.RFC822Z,
		"RFC850":                        time.RFC850,
		"RFC1123":                       time.RFC1123,
		"RFC1123Z":                      time.RFC1123Z,
		"Kitchen":                       time.Kitchen,
		"Stamp":                         time.Stamp,
		"StampMilli":                    time.StampMilli,
		"StampMicro":                    time.StampMicro,
		"StampNano":                     time.StampNano,
		"DateTime":                      time.DateTime,
		"Unix":                          unixStrategy,
		"UnixMicro":                     unixStrategy,
		"UnixMilli":                     unixStrategy,
	}
)

var ErrUnknownFormat = errors.New("unknown format")

// formatTime turns a time string into a prettier time string.
//
// Args:
//   - timeStr: the incoming time string
//   - inputFormat: the format of the incoming time string
//   - outputFormat: the pretty string format. More info at the
//     [go time format documentation](https://go.dev/src/time/format.go)
//
// Returns:
// - a string
//
// If parsing did not succeed or the inputFormat is unknown, returns the value of timeStr.
//
// Output format convention by go uses numbers instead of strings like H, m or YYYY
// Hour                                       "15"
// ZeroHour12                                 "03"
// ZeroMinute                                 "04"
// ZeroSecond                                 "05"
// Fractional Seconds (incl. trailing zeros)  ".00" (any amount of digits; up to 9).
func formatTime(timeStr, inputFormat, outputFormat string) string {
	var (
		parsedTime time.Time
		err        error
	)

	layout, ok := formats[inputFormat]
	if !ok {
		layout = inputFormat
	}

	switch layout {
	case unixStrategy:
		parsedTime, err = parseUnix(inputFormat, timeStr)
	default:
		parsedTime, err = time.Parse(layout, timeStr)
	}

	if err != nil {
		return timeStr
	}

	return parsedTime.Format(outputFormat)
}

// parseUnix takes a timestamp string and parses it.
func parseUnix(format, timestamp string) (time.Time, error) {
	switch format {
	case "Unix":
		split := strings.Split(timestamp, ".")
		if len(split) > 2 {
			return time.Time{}, ErrUnknownFormat
		}

		sec, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("can not parse sec: %w", err)
		}

		var nsec int64 = 0

		if len(split) == 2 {
			nsecString := toNanoSec(split[1])

			nsec, err = strconv.ParseInt(nsecString, 10, 64)
			if err != nil {
				return time.Time{}, fmt.Errorf("can not parse nsecString: %w", err)
			}
		}

		return time.Unix(sec, nsec), nil

	case "UnixMilli":
		i64, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("can not parse timestamp: %w", err)
		}

		return time.UnixMilli(i64), nil

	case "UnixMicro":
		i64, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("can not parse timestamp: %w", err)
		}

		return time.UnixMicro(i64), nil
	}

	return time.Time{}, ErrUnknownFormat
}

func toNanoSec(value string) string {
	return rightPad(value, 9)
}

// rightPad pads a string with x-`count` 0 (zeroes) to the right
//
// If count is smaller than len(value), the value will be cut off from the
// right.
func rightPad(value string, count int) string {
	if count < 0 {
		count = 0
	}

	if len(value) == count {
		return value
	}

	if len(value) > count {
		return value[:count]
	}

	paddingCount := count - len(value)
	padding := strings.Repeat("0", paddingCount)

	return value + padding
}
