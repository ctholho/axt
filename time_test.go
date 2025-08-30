package main

import (
	"testing"
	"time"
)

func TestRightPad(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		inputString string
		inputCount  int
		expected    string
	}{
		{
			name:        "Basic Padding",
			inputString: "1234",
			inputCount:  10,
			expected:    "1234000000",
		},
		{
			name:        "Basic Padding 2",
			inputString: "1234",
			inputCount:  5,
			expected:    "12340",
		},
		{
			name:        "Cut off padding",
			inputString: "12345678",
			inputCount:  5,
			expected:    "12345",
		},
		{
			name:        "Exact Length",
			inputString: "12345",
			inputCount:  5,
			expected:    "12345",
		},
		{
			name:        "Empty string",
			inputString: "",
			inputCount:  5,
			expected:    "00000",
		},
		{
			name:        "Zero Count",
			inputString: "1245",
			inputCount:  0,
			expected:    "",
		},
		{
			name:        "Negative count is like zero count",
			inputString: "test",
			inputCount:  -5,
			expected:    "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual := rightPad(testCase.inputString, testCase.inputCount)
			if actual != testCase.expected {
				t.Errorf("rightPad(\"%s\", %d) = \"%s\"; want \"%s\"", testCase.inputString, testCase.inputCount, actual, testCase.expected)
			}
		})
	}
}

func TestParseUnix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		timeStr  string
		format   string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "Unix format - seconds only",
			timeStr:  "1756555555",
			format:   "Unix",
			wantTime: time.Unix(1756555555, 0),
			wantErr:  false,
		},
		{
			name:     "Unix format - with nanoseconds",
			timeStr:  "1756555555.123456789",
			format:   "Unix",
			wantTime: time.Unix(1756555555, 123456789),
			wantErr:  false,
		},
		{
			name:     "Unix format - with partial nanoseconds (needs padding)",
			timeStr:  "1756555555.123",
			format:   "Unix",
			wantTime: time.Unix(1756555555, 123000000),
			wantErr:  false,
		},
		{
			name:     "UnixMilli format - valid",
			timeStr:  "1756555555123",
			format:   "UnixMilli",
			wantTime: time.UnixMilli(1756555555123),
			wantErr:  false,
		},
		{
			name:     "UnixMicro format - valid",
			timeStr:  "1756555555123456",
			format:   "UnixMicro",
			wantTime: time.UnixMicro(1756555555123456),
			wantErr:  false,
		},
		{
			name:     "Error - Unknown format",
			timeStr:  "1756555555",
			format:   "blabla",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "Error - Unix format with invalid seconds",
			timeStr:  "not-a-number.123",
			format:   "Unix",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "Error - Unix format with invalid nanoseconds",
			timeStr:  "1756555555.not-a-number",
			format:   "Unix",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "Error - Unix format with too many separators",
			timeStr:  "1756555555.123.456",
			format:   "Unix",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "Error - UnixMilli format with non-integer string",
			timeStr:  "not-a-number",
			format:   "UnixMilli",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "Error - UnixMicro format with non-integer string",
			timeStr:  "not-a-number",
			format:   "UnixMicro",
			wantTime: time.Time{},
			wantErr:  true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			gotTime, err := parseUnix(testCase.format, testCase.timeStr)
			if (err != nil) != testCase.wantErr {
				t.Errorf("parseUnix() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !testCase.wantErr && !gotTime.Equal(testCase.wantTime) {
				t.Errorf("parseUnix() gotTime = %v, want %v", gotTime, testCase.wantTime)
			}
		})
	}
}
