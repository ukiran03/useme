package menu

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
)

type FileItem struct {
	Name     string
	Selected bool
}

type FileSet struct {
	files []*FileItem
}

func (f *FileItem) String() string {
	var status string
	switch f.Selected {
	case true:
		status = "#"
	case false:
		status = ""

	}
	return f.Name + status
}

func NewFileSet(files []*FileItem) *FileSet {
	return &FileSet{
		files: files,
	}
}

func main() {
	files := []*FileItem{
		{"main.go", false},
		{"utils.go", false},
		{"README.md", false},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		// The "View" (Print the menu)
		printMenu(files)

		// The "Read" (Get user input)
		fmt.Printf("%s\nWhat now> %s", colors[Blue], colors[Reset])
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		// The "Eval" (Handle logic)
		switch input {
		case "q", "quit":
			fmt.Println("Bye!")
			return
		case "u", "update":
			fmt.Println()
		default:
			// If input is a number, toggle that file
			handleSelection(input, files)
		}
	}
}

func printMenu(files []*FileItem) {
	fmt.Printf("%s\n*** Interactive File Picker ***%s\n",
		colors[BoldText], colors[Reset])

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for i, f := range files {
		status := " "
		pre, post := "", ""
		if f.Selected {
			status = "*"
			pre = colors[BoldText] + colors[Blue]
			post = colors[Reset]
		}
		label := fmt.Sprintf("%s%s %d: %s%s", pre, status, i+1, f.Name, post)
		fmt.Fprintf(w, "\t%s\n", label)
	}

	w.Flush()
	fmt.Printf("\n%s(Commands: 1-%d to toggle, 'u' update, 'q' quit): %s",
		colors[BoldText], len(files), colors[Reset])
}

func handleSelection(input string, files []*FileItem) {
	indexes := parseOnlyInts(input)
	for _, idx := range indexes {
		if idx > 0 && idx <= len(files) {
			files[idx-1].Selected = !files[idx-1].Selected // toggle
		} else {
			fmt.Printf("Unknown command: %s\n", input)
		}
	}
}

func parseOnlyInts(in string) []int {
	// This finds all sequences of one or more digits
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(in, -1)

	results := make([]int, 0, len(matches))
	for _, m := range matches {
		n, _ := strconv.Atoi(m)
		results = append(results, n)
	}
	return results
}
