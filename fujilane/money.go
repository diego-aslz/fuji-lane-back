package fujilane

import humanize "github.com/dustin/go-humanize"

// FormatCents returns the given cents formatted as money
func FormatCents(cents int) string {
	return "$" + humanize.FormatFloat("#,###.##", float64(cents)/100.0)
}
