package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// apiError models exactly the Foxbit error payload:
type apiError struct {
	Error struct {
		Message string   `json:"message"`
		Code    int      `json:"code"`
		Details []string `json:"details"`
	} `json:"error"`
}

// DisplayError prints any Foxbit‐style error in a single table:
// TYPE | CODE | MESSAGE | DETAIL
func DisplayError(err error) {
	raw := err.Error()
	// find the JSON blob
	idx := strings.Index(raw, "{")
	if idx < 0 {
		fmt.Fprintln(os.Stderr, raw)
		return
	}

	blob := raw[idx:]
	var e apiError
	if jsonErr := json.Unmarshal([]byte(blob), &e); jsonErr != nil {
		// malformed JSON: fallback to raw
		fmt.Fprintln(os.Stderr, raw)
		return
	}

	w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
	// header
	fmt.Fprintln(w, "TYPE\tCODE\tMESSAGE\tDETAIL")
	// one row per detail (or a single empty‐detail row)
	if len(e.Error.Details) > 0 {
		for _, d := range e.Error.Details {
			fmt.Fprintf(w, "ERROR\t%d\t%s\t%s\n",
				e.Error.Code,
				e.Error.Message,
				d,
			)
		}
	} else {
		fmt.Fprintf(w, "ERROR\t%d\t%s\t\n",
			e.Error.Code,
			e.Error.Message,
		)
	}
	w.Flush()
}
