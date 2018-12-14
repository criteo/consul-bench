package bench

import (
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/consul/agent/pool"
	"github.com/hashicorp/consul/agent/structs"
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

func QueryServer(addr string, dc string, serviceName string, wait time.Duration, allowStale bool) queryFn {
	connPool := &pool.ConnPool{
		SrcAddr:    nil,
		LogOutput:  os.Stderr,
		MaxTime:    time.Hour,
		MaxStreams: 1000000,
		TLSWrapper: nil,
		ForceTLS:   false,
	}

	args := structs.ServiceSpecificRequest{
		Datacenter:  dc,
		ServiceName: serviceName,
		Source: structs.QuerySource{
			Datacenter: dc,
			Node:       "test-1",
		},
		QueryOptions: structs.QueryOptions{
			MaxQueryTime: 10 * time.Minute,
		},
	}

	ip, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal(err)
	}
	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 8300
	}
	srvAddr := &net.TCPAddr{net.ParseIP(ip), port, ""}

	return func(index uint64) (uint64, error) {
		args.QueryOptions.MinQueryIndex = index
		var resp *structs.IndexedCheckServiceNodes
		err := connPool.RPC(dc, srvAddr, 3, "Health.ServiceNodes", false, &args, &resp)
		if err != nil {
			return 0, err
		}
		return resp.QueryMeta.Index, nil
	}
}
