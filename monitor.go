package main

import (
	"log"
	"time"

	"github.com/criteo/consul-bench/stats"
	"github.com/shirou/gopsutil/process"
)

func Monitor(pid int32, statsC chan stats.Stat, done chan struct{}) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		log.Fatal(err)
	}
	proc.Percent(0)

	tick := time.Tick(time.Second)
	for {
		select {
		case <-done:
			return
		case <-tick:
			p, err := proc.Percent(0)
			if err != nil {
				log.Println(err)
			} else {
				select {
				case statsC <- stats.Stat{
					Label: "CPU",
					Value: p,
				}:
				case <-done:
					return
				}
			}
		}
	}
}
