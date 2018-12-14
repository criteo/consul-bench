package bench

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

type Stat struct {
	Label string
	Value float64
}

func DisplayStats(ch <-chan Stat) {
	var l sync.Mutex
	entries := map[string]Stat{}

	go func() {
		for {
			e := <-ch
			l.Lock()
			entries[e.Label] = e
			l.Unlock()
		}
	}()

	for range time.Tick(time.Second) {
		l.Lock()
		var entriesSlice []Stat
		for _, e := range entries {
			entriesSlice = append(entriesSlice, e)
		}

		sort.Slice(entriesSlice, func(i, j int) bool {
			return strings.Compare(entriesSlice[i].Label, entriesSlice[j].Label) < 0
		})
		var s []string
		for _, e := range entriesSlice {
			s = append(s, fmt.Sprintf("%s: %f", e.Label, e.Value))
		}

		log.Println(strings.Join(s, ", "))
		entriesSlice = entriesSlice[:0]
		l.Unlock()
	}
}
