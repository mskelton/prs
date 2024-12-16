package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
)

type CellColor string

const (
	CellColorDefault CellColor = ""
	CellColorRed     CellColor = "red"
	CellColorGreen   CellColor = "green"
	CellColorYellow  CellColor = "yellow"
	CellColorBlue    CellColor = "blue"
	CellColorMagenta CellColor = "magenta"
	CellColorCyan    CellColor = "cyan"
	CellColorGray    CellColor = "gray"
	CellColorDim     CellColor = "dim"
)

func getColor(c CellColor) *color.Color {
	switch c {
	case CellColorRed:
		return color.New(color.FgRed)
	case CellColorGreen:
		return color.New(color.FgGreen)
	case CellColorYellow:
		return color.New(color.FgYellow)
	case CellColorBlue:
		return color.New(color.FgBlue)
	case CellColorMagenta:
		return color.New(color.FgMagenta)
	case CellColorCyan:
		return color.New(color.FgCyan)
	case CellColorGray:
		return color.RGB(99, 101, 123)
	case CellColorDim:
		return color.New(color.FgWhite)
	}

	return color.New(color.FgHiWhite)
}

const separator = "  "

type Cell struct {
	Value string
	Color CellColor
}

type Row struct {
	Cells []Cell
}

type Table struct {
	Columns []string
	Rows    []Row
}

// Special implementation of string padding to account for unicode string width
func pad(str string, w int) string {
	return str + strings.Repeat(" ", w-runewidth.StringWidth(str))
}

func toUpper(str string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	withSpaces := re.ReplaceAllString(str, "${1} ${2}")
	return strings.ToUpper(withSpaces)
}

func (table *Table) Print() {
	widths := make([]int, len(table.Columns))
	headerColor := getColor(CellColorGray).Add(color.Underline).SprintFunc()

	// Find the maximum width of each column
	for _, row := range table.Rows {
		for i, cell := range row.Cells {
			length := runewidth.StringWidth(cell.Value)
			widths[i] = max(widths[i], length)
		}
	}

	// Calculate the width of each column header, ignoring empty columns
	for i, col := range table.Columns {
		if widths[i] > 0 {
			// Column headers never have Unicode, so `len()` is safe to use
			widths[i] = max(widths[i], len(col))
		}
	}

	// Create the header row, skipping empty columns
	var header []string
	for i, col := range table.Columns {
		if widths[i] > 0 {
			header = append(header, headerColor(pad(toUpper(col), widths[i])))
		}
	}

	fmt.Println(strings.Join(header, "  "))

	// Print an ASCII underline if colorization is disabled
	if color.NoColor {
		var underline []string

		for _, width := range widths {
			if width > 0 {
				underline = append(underline, strings.Repeat("-", width))
			}
		}

		fmt.Println(strings.Join(underline, "  "))
	}

	for _, row := range table.Rows {
		var cells []string

		for i, cell := range row.Cells {
			if widths[i] > 0 {
				color := getColor(cell.Color).SprintFunc()

				cells = append(cells, color(pad(cell.Value, widths[i])))
			}
		}

		fmt.Println(strings.Join(cells, "  "))
	}
}
