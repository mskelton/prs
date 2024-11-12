package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type PR struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Branch    string `json:"headRefName"`
	CreatedAt string `json:"createdAt"`
	State     string `json:"state"`
	Url       string `json:"url"`
}

func timeAgo(createdAt string) string {
	t, _ := time.Parse(time.RFC3339, createdAt)
	duration := time.Since(t)
	hours := int(duration.Hours())

	if hours == 0 {
		return fmt.Sprintf("about %d minutes ago", int(duration.Minutes()))
	}

	return fmt.Sprintf("about %d hours ago", hours)
}

func createTableRow(pr PR) Row {
	var numberColor CellColor
	if pr.State == "OPEN" {
		numberColor = CellColorGreen
	} else {
		numberColor = CellColorMagenta
	}

	return Row{
		Cells: []Cell{
			{Value: fmt.Sprintf("#%d", pr.Number), Color: numberColor},
			{Value: pr.Title},
			{Value: pr.Branch, Color: CellColorCyan},
			{Value: pr.Url, Color: CellColorBlue},
			{Value: timeAgo(pr.CreatedAt)},
		},
	}
}

func main() {
	// Fetch PR data from GitHub CLI
	cmd := exec.Command("gh", "pr", "list", "--json", "number,title,headRefName,createdAt,state,url")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Parse JSON data into PR struct
	var prs []PR
	if err := json.Unmarshal(output, &prs); err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	table := Table{
		Columns: []string{"ID", "TITLE", "BRANCH", "URL", "CREATED AT"},
		Rows:    []Row{},
	}

	for _, pr := range prs {
		table.Rows = append(table.Rows, createTableRow(pr))
	}

	table.Print()
}
