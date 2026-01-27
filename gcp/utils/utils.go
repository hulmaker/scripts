package utils

import (
	"path"
	"net/url"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)


// clamps text to given width
func Clamp(text string, width uint) string {
	text = strings.Trim(text, " ")
	w := int(width) // w >= 0
	if len(text) <= w {
		return text + strings.Repeat(" ", w-len(text))
	}
	return text[:w]
}

// CellFormatter Expects: column_name, text. (text is already padded with spaces to the max column width)
type CellFormatter func(string, string) string

// Print the data as a table in human readable format
// Columns are separated by the separator
// cell_formatter to override the default text formatting
// NOTE: text/tabwriter messes up the column alignment when ANSI colors are used.
func PrintTable(data [][]string, separator string, formatter CellFormatter) error {
	if len(data) == 0 {
		return nil
	}

	// calculate the column columnWidths
	columnWidths := make([]uint, len(data[0]))
	for _, row := range data {
		for j, cell := range row {
			columnWidths[j] = max(columnWidths[j], uint(len(cell)))
		}
	}

	// print the table
	for _, row := range data {
		arr := make([]string, len(row))
		for j, text := range row {
			arr[j] = Clamp(text, columnWidths[j])
			if formatter != nil {
				arr[j] = formatter(data[0][j], arr[j])
			}
		}
		fmt.Println(strings.Join(arr, separator))
	}
	return nil
}


// A safe way to create escaped version of url: <part_0>/<part_1>/.../<parts_n>/?<key_1>=<val_1>&<key_2>=<val_2>
func BuildUrl(parts []string, args map[string]string) (string, error) {
	u, err := url.Parse(path.Join(parts...))
	if err != nil {
		return "", err
	}
	query := u.Query()
	for key, val := range args {
		_key := strings.TrimSpace(key)
		_val := strings.TrimSpace(val)
		if _key == "" || _val == "" {
			continue
		}
		query.Set(_key, _val)
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}


func OpenBrowser(url string) error {
	fmt.Printf("Opening url in default browser. If nothing happens, plase visit:\n%s\n", url)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}

func RunCommands(commands [][]string) error {
	for _, cmdArgs := range commands {
		if len(cmdArgs) == 0 {
			return fmt.Errorf("can't execute an empty command")
		}
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
