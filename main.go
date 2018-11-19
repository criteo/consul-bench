package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	consul "github.com/hashicorp/consul/api"
)

const (
	InstanceCount = 800
	Watchers      = 150
	FlapInterval  = 800 * time.Second
)

type IdxStat struct {
	Acks      int
	ChangedAt time.Time
	FirstAck  time.Time
	LastAck   time.Time
}

func main() {
	c, err := consul.NewClient(&consul.Config{})
	if err != nil {
		log.Fatal(err)
	}

	idx := uint64(0)
	stats := map[uint64]IdxStat{}
	statsLock := sync.Mutex{}

	for instanceID := 0; instanceID < InstanceCount; instanceID++ {
		err := c.Agent().ServiceRegister(&consul.AgentServiceRegistration{
			Name: "srv",
			ID:   fmt.Sprintf("srv-%d", instanceID),
			Checks: []*consul.AgentServiceCheck{
				{
					CheckID: fmt.Sprintf("check-%d", instanceID),
					TTL:     "1m",
				},
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	for instanceID := 0; instanceID < InstanceCount; instanceID++ {
		go func(instanceID int) {
			time.Sleep(FlapInterval / time.Duration(InstanceCount) * time.Duration(instanceID))
			lastStatus := false
			for {
				nidx := atomic.AddUint64(&idx, 1)
				var f func(checkID, note string) error
				if lastStatus {
					f = c.Agent().FailTTL
				} else {
					f = c.Agent().PassTTL
				}

				err := f(fmt.Sprintf("check-%d", instanceID), fmt.Sprint(nidx))
				if err != nil {
					log.Fatal(err)
				}

				statsLock.Lock()
				stats[nidx] = IdxStat{
					ChangedAt: time.Now(),
				}
				statsLock.Unlock()

				lastStatus = !lastStatus

				time.Sleep(FlapInterval)
			}
		}(instanceID)
	}

	log.Println("services registered")

	for watcherID := 0; watcherID < Watchers; watcherID++ {
		go func() {
			lastIdx := uint64(0)

			for {
				entries, meta, err := c.Health().Service("srv", "", false, &consul.QueryOptions{
					WaitTime:  10 * time.Minute,
					WaitIndex: lastIdx,
				})
				if err != nil {
					log.Fatal(err)
				}

				lastIdx = meta.LastIndex

				//maxIdx := uint64(0)
				for _, entry := range entries {
					for _, c := range entry.Checks {
						if c.ServiceID == "" {
							continue
						}

						if c.Notes == "" {
							continue
						}

						fmt.Println("Check: ", c.Notes)
					}
				}
			}
		}()
	}

	<-make(chan bool)
}
