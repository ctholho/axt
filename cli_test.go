package main

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// ansiRegex strips ANSI color codes from the output.
func stripAnsi(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// captureOutput runs the main function with specified arguments and input,
// and returns the captured standard output.
func captureOutput(t *testing.T, args []string, input string) string {
	t.Helper()
	// Keep track of the original os variables to restore them later
	oldArgs := os.Args
	oldStdin := os.Stdin
	oldStdout := os.Stdout

	defer func() {
		os.Args = oldArgs
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	rIn, wIn, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}

	os.Stdin = rIn

	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdout = wOut

	os.Args = append([]string{"axt"}, args...)

	go func() {
		defer wIn.Close()

		_, err := wIn.WriteString(input)
		if err != nil {
			t.Errorf("failed to write to stdin pipe: %v", err)
		}
	}()

	main()

	wOut.Close()

	var buf bytes.Buffer

	_, err = io.Copy(&buf, rOut)
	if err != nil {
		t.Errorf("failed to read from stdout pipe: %v", err)
	}

	return buf.String()
}

//nolint:paralleltest // meddles with system globals and can't run in parallel
func TestMainFunction(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		input    string
		expected string // Expected output, can be multi-line.
		useUTC   bool
	}{
		{
			name:   "Default slog format",
			args:   []string{},
			useUTC: false,
			input:  `{"time":"2025-08-24T21:51:45.549605+02:00","level":"INFO","msg":"API request completed","status_code":200,"response_time_ms":127}`,
			expected: `
 21:51:45.549 INFO  API request completed
                   status_code: 200
                   response_time_ms: 127
`,
		},
		{
			name:     "Unstructured non-JSON log",
			args:     []string{},
			useUTC:   false,
			input:    `something without proper JSON`,
			expected: `🪵  something without proper JSON`,
		},
		{
			name:   "Custom flags for ECS",
			args:   []string{"-t", "@timestamp", "-l", "log.level", "-m", "message"},
			useUTC: false,
			input:  `{"@timestamp":"2025-08-24T21:51:45.549Z","log.level":"ERROR","message":"User authentication failed","error.message":"invalid credentials"}`,
			expected: `
 21:51:45.549 ERR   User authentication failed
                   error.message: "invalid credentials"
`,
		},
		{
			name:   "Emoji flag for levels",
			args:   []string{"--emoji"},
			useUTC: false,
			input:  `{"time":"2025-08-24T21:51:45.549Z","level":"WARN","msg":"Deprecated API used"}`,
			expected: `
 21:51:45.549 ⚠️  Deprecated API used
`,
		},
		{
			name:   "Hide properties from output",
			args:   []string{"--hide", "trace_id", "--hide", "user_agent"},
			useUTC: false,
			input:  `{"time":"2025-08-24T21:51:45.549Z","level":"INFO","msg":"Request received","trace_id":"xyz","user_agent":"test-runner","important": true}`,
			expected: `
 21:51:45.549 INFO  Request received
                   important: true
`,
		},
		{
			name:   "Custom time out format",
			args:   []string{"--time-out", "2006/01/02 15h04m05s.000"},
			useUTC: true,
			input:  `{"time":"2025-08-24T21:51:45.549Z","level":"DEBUG","msg":"whatever floats your boat time format test"}`,
			expected: `
 2025/08/24 21h51m45s.549 DEBUG whatever floats your boat time format test
`,
		},
		{
			name:   "Custom time in format - UnixMilli",
			args:   []string{"--time-in", "UnixMilli", "--time-out", "15:04:05.000"},
			useUTC: true,
			input:  `{"time":"1756555555123","level":"DEBUG","msg":"Timestamp test"}`,
			expected: `
		  12:05:55.123 DEBUG Timestamp test
		 `,
		},
		{
			name:   "Custom time in format - Unix with decimal",
			args:   []string{"--time-in", "Unix"},
			useUTC: true,
			input:  `{"time":"1756555555.123","level":"DEBUG","msg":"Timestamp test"}`,
			expected: `
		  12:05:55.123 DEBUG Timestamp test
		 `,
		},
	}

	for _, testCase := range testCases {
		//nolint:paralleltest // meddles with system globals and can't run in parallel
		t.Run(testCase.name, func(t *testing.T) {
			// Reset the command-line flags before each test run.
			// The main() function registers flags, and calling it in a loop
			// would cause a "flag redefined" panic if we didn't reset.
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

			actualRaw := captureOutput(t, testCase.args, testCase.input)
			actualClean := strings.TrimSpace(stripAnsi(actualRaw))
			expectedClean := strings.TrimSpace(testCase.expected)

			expectedLines := strings.Split(expectedClean, "\n")
			actualLines := strings.Split(actualClean, "\n")

			// Compare the first line (the header) directly
			if strings.TrimSpace(actualLines[0]) != strings.TrimSpace(expectedLines[0]) {
				t.Errorf("Header line mismatch.\n--- Expected ---\n%s\n--- Actual ---\n%s", expectedLines[0], actualLines[0])

				return
			}

			// For the remaining lines (properties), their order isn't guaranteed.
			// We just check that each expected property line is present in the actual output.
			if len(expectedLines) > 1 {
				for _, expectedLine := range expectedLines[1:] {
					trimmedExpected := strings.TrimSpace(expectedLine)
					found := false

					for _, actualLine := range actualLines[1:] {
						if strings.TrimSpace(actualLine) == trimmedExpected {
							found = true

							break
						}
					}

					if !found {
						t.Errorf("Expected property line not found in output.\n--- Missing Line ---\n%s\n--- Actual Output ---\n%s", trimmedExpected, actualClean)
					}
				}
			}
		})
	}
}
