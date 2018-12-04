package main

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	consul "github.com/hashicorp/consul/api"
)

type queryFn func(uint64) (uint64, error)

func RunQueries(fn queryFn, count int, stats chan Stat) error {
	log.Println("Starting", count, "watchers...")

	var qps int32

	errs := make(chan error, 1)
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-done:
				return
			default:
			}

			index := uint64(0)
			var err error
			for {
				index, err = fn(index)
				if err != nil {
					select {
					case errs <- err:
					default:
					}
					return
				}
				atomic.AddInt32(&qps, 1)
			}
		}()
	}
	go func() {
		for range time.Tick(time.Second) {
			c := atomic.SwapInt32(&qps, 0)
			stats <- Stat{"QPS", float64(c)}
		}
	}()
	log.Println("Watchers started.")
	wg.Wait()
	select {
	case err := <-errs:
		return err
	default:
	}
	return nil
}

func QueryAgent(client *consul.Client, serviceName string, wait time.Duration, allowStale bool) queryFn {
	return func(index uint64) (uint64, error) {
		_, meta, err := client.Health().Service(serviceName, "", false, &consul.QueryOptions{
			WaitTime:   wait,
			WaitIndex:  index,
			AllowStale: allowStale,
		})
		if err != nil {
			return 0, err
		}

		return meta.LastIndex, nil
	}
}
