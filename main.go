package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mergestat/timediff"
)

type PullRequest struct {
	Branch    string `json:"headRefName"`
	CreatedAt string `json:"createdAt"`
	IsDraft   bool   `json:"isDraft"`
	Number    int    `json:"number"`
	State     string `json:"state"`
	Title     string `json:"title"`
	Url       string `json:"url"`
}

func timeAgo(createdAt string) string {
	t, _ := time.Parse(time.RFC3339, createdAt)

	return timediff.TimeDiff(t)
}

func colorForPRState(pr PullRequest) CellColor {
	switch pr.State {
	case "OPEN":
		if pr.IsDraft {
			return "gray"
		}
		return "green"
	case "CLOSED":
		return "red"
	case "MERGED":
		return "magenta"
	default:
		return ""
	}
}

func createTableRow(pr PullRequest) Row {
	return Row{
		Cells: []Cell{
			{Value: fmt.Sprintf("#%d", pr.Number), Color: colorForPRState(pr)},
			{Value: pr.Title},
			{Value: pr.Url, Color: CellColorBlue},
			{Value: timeAgo(pr.CreatedAt), Color: CellColorDim},
		},
	}
}

func main() {
	cmd := exec.Command("gh", "pr", "list", "--json", "number,title,headRefName,createdAt,state,url,isDraft")
	cmd.Args = append(cmd.Args, os.Args[1:]...)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		os.Exit(1)
	}

	// Parse JSON data into PR struct
	var prs []PullRequest
	if err := json.Unmarshal(output, &prs); err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	table := Table{
		Columns: []string{"ID", "TITLE", "URL", "CREATED AT"},
		Rows:    []Row{},
	}

	for _, pr := range prs {
		table.Rows = append(table.Rows, createTableRow(pr))
	}

	table.Print()
}
