package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mergestat/timediff"
)

type PullRequest struct {
	Branch    string `json:"headRefName"`
	CreatedAt string `json:"createdAt"`
	ClosedAt  string `json:"closedAt"`
	IsDraft   bool   `json:"isDraft"`
	Number    int    `json:"number"`
	State     string `json:"state"`
	Title     string `json:"title"`
	Url       string `json:"url"`
}

func timeAgo(createdAt string) string {
	if createdAt == "" {
		return ""
	}

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

func createTableRow(pr PullRequest, columns []string) Row {
	cells := []Cell{}

	for _, column := range columns {
		switch column {
		case "id":
			cells = append(cells, Cell{Value: fmt.Sprintf("#%d", pr.Number), Color: colorForPRState(pr)})
		case "title":
			cells = append(cells, Cell{Value: pr.Title})
		case "url":
			cells = append(cells, Cell{Value: pr.Url, Color: CellColorBlue})
		case "createdAt":
			cells = append(cells, Cell{Value: timeAgo(pr.CreatedAt), Color: CellColorDim})
		case "closedAt":
			cells = append(cells, Cell{Value: timeAgo(pr.ClosedAt), Color: CellColorDim})
		}
	}

	return Row{Cells: cells}
}

func main() {
	cmd := exec.Command("gh", "pr", "list", "--json", "number,title,headRefName,createdAt,closedAt,state,url,isDraft")
	cmd.Stderr = os.Stderr

	columns := []string{"id", "title", "url", "createdAt", "closedAt"}

	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "--columns=") {
		columns = strings.Split(strings.TrimPrefix(os.Args[1], "--columns="), ",")
		cmd.Args = append(cmd.Args, os.Args[2:]...)
	} else {
		cmd.Args = append(cmd.Args, os.Args[1:]...)
	}

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
		Columns: columns,
		Rows:    []Row{},
	}

	for _, pr := range prs {
		table.Rows = append(table.Rows, createTableRow(pr, columns))
	}

	table.Print()
}
