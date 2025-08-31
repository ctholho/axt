package main

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

// addBorder shows a pretty border around logged properties.
func addBorder(logLines []string, verticalLine string) {
	if len(logLines) > 3 {
		outerColor := pterm.NewRGB(70, 70, 70)
		innerColor := pterm.NewRGB(150, 150, 150)

		// Fade for the line
		for index, line := range logLines {
			var currentVerticalLine string

			switch index {
			case 0:
				currentVerticalLine = outerColor.Fade(0, float32(len(logLines)-1), float32(index), innerColor, outerColor).Sprint("     ┌")
			case len(logLines) - 1:
				currentVerticalLine = outerColor.Fade(0, float32(len(logLines)-1), float32(index), innerColor, outerColor).Sprint("     └")
			default:
				currentVerticalLine = outerColor.Fade(0, float32(len(logLines)-1), float32(index), innerColor, outerColor).Sprint("     │")
			}

			fmt.Println(strings.Replace(line, verticalLine, currentVerticalLine, 1))
		}
	} else {
		for _, line := range logLines {
			fmt.Println(line)
		}
	}
}
