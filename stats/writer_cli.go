package stats

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

type CLIWriter struct {
}

func (*CLIWriter) Write(entries []Stat) error {
	sort.Slice(entries, func(i, j int) bool {
		return strings.Compare(entries[i].Label, entries[j].Label) < 0
	})
	var s []string
	for _, e := range entries {
		s = append(s, fmt.Sprintf("%s: %.2f", e.Label, e.Value))
	}

	log.Println(strings.Join(s, ", "))

	return nil
}
