package stats

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

type statC struct {
	Stat
	Count int
}

type Stats struct {
	avg     map[string]*statC
	entries map[string]Stat
	l       sync.Mutex
}

func New() *Stats {
	return &Stats{
		avg:     make(map[string]*statC),
		entries: make(map[string]Stat),
	}
}

func (s *Stats) Run(ch chan Stat) {
	done := make(chan struct{})

	go func() {
		for e := range ch {
			s.l.Lock()
			s.entries[e.Label] = e
			if _, ok := s.avg[e.Label]; !ok {
				s.avg[e.Label] = &statC{
					Count: 0,
					Stat:  e,
				}
			}
			s.avg[e.Label].Value = (s.avg[e.Label].Value*float64(s.avg[e.Label].Count) + e.Value) / float64(s.avg[e.Label].Count+1)
			s.avg[e.Label].Count++
			s.l.Unlock()
		}
		close(done)
	}()

	start := time.Now()
	tick := time.Tick(time.Second)

	for {
		select {
		case <-done:
			s.l.Lock()
			entriesSlice := []Stat{{
				Label: "Runtime (s)",
				Value: time.Since(start).Seconds(),
			}}
			for _, e := range s.avg {
				entriesSlice = append(entriesSlice, e.Stat)
			}

			log.Println("====== Summary ======")
			printLine(entriesSlice)
			s.l.Unlock()
			return
		case <-tick:

			s.l.Lock()
			var entriesSlice []Stat
			for _, e := range s.entries {
				entriesSlice = append(entriesSlice, e)
			}

			printLine(entriesSlice)

			entriesSlice = entriesSlice[:0]
			s.l.Unlock()
		}
	}
}

func (s *Stats) AVGs() []Stat {
	s.l.Lock()
	defer s.l.Unlock()

	entries := []Stat{}
	for _, v := range s.avg {
		entries = append(entries, v.Stat)
	}

	sort.Slice(entries, func(i, j int) bool {
		return strings.Compare(entries[i].Label, entries[j].Label) < 0
	})

	return entries
}

func (s *Stats) Reset() {
	s.l.Lock()
	defer s.l.Unlock()
	s.avg = make(map[string]*statC)
	s.entries = make(map[string]Stat)
}

func printLine(entries []Stat) {
	sort.Slice(entries, func(i, j int) bool {
		return strings.Compare(entries[i].Label, entries[j].Label) < 0
	})
	var s []string
	for _, e := range entries {
		s = append(s, fmt.Sprintf("%s: %.2f", e.Label, e.Value))
	}

	log.Println(strings.Join(s, ", "))
}
