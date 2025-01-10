package overlay

// Take from https://gist.github.com/Broderick-Westrope/b89b14770c09dda928c4a108f437b927

// Write own simplified implementation
// Read the code, see how it works, and do it from memory. write tests.

// import (
// 	"fmt"
// 	"regexp"
// 	"strings"
//
// 	"github.com/charmbracelet/lipgloss"
// 	"github.com/charmbracelet/x/ansi"
// )
//
// // Overlays a foreground over a background, aligned center
// func Overlay(bg string) {
// 	fg := "test"
//
// 	bgLines := strings.Split(bg, "\n")
// 	fgLines := strings.Split(fg, "\n")
//
// 	for i, fgLine := range fgLines {
// 		bgLine := bgLines[i+4] // 4 = row
//
//         if len(bgLine) < 20 { // 20 = col
//             bgLine += strings.Repeat()
//         }
// 	}
// }
