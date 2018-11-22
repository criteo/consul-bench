package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	consul "github.com/hashicorp/consul/api"
)

func main() {
	consulAddr := flag.String("consul", "127.0.0.1:8500", "Consul address")
	serviceName := flag.String("service", "srv", "Service to watch")
	registerInstances := flag.Int("register", 0, "Register N -service instances")
	flapInterval := flag.Duration("flap-interval", 0, "If -register is given, flap each instance between critical and passing state on given interval")
	watchers := flag.Int("watchers", 1, "Number of concurrnet watchers on service")
	flag.Parse()

	c, err := consul.NewClient(&consul.Config{
		Address: *consulAddr,
	})
	if err != nil {
		log.Fatal(err)
	}

	startCh := make(chan struct{})

	if *registerInstances > 0 {
		log.Printf("Registering %d %s instances...\n", *registerInstances, *serviceName)

		for instanceID := 0; instanceID < *registerInstances; instanceID++ {
			err := c.Agent().ServiceRegister(&consul.AgentServiceRegistration{
				Name: "srv",
				ID:   fmt.Sprintf("srv-%d", instanceID),
				Checks: []*consul.AgentServiceCheck{
					{
						CheckID: fmt.Sprintf("check-%d", instanceID),
						TTL:     "1m",
						Status:  consul.HealthCritical,
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}

		if *flapInterval > 0 {
			log.Printf("Flapping instances every %s", *flapInterval)

			for instanceID := 0; instanceID < *registerInstances; instanceID++ {
				go func(instanceID int) {
					<-startCh

					time.Sleep((*flapInterval / time.Duration(*registerInstances)) * time.Duration(instanceID))
					lastStatus := false
					for {
						var f func(checkID, note string) error
						if lastStatus {
							f = c.Agent().FailTTL
							fmt.Print("F")
						} else {
							f = c.Agent().PassTTL
							fmt.Print("P")
						}

						err := f(fmt.Sprintf("check-%d", instanceID), "")
						if err != nil {
							log.Fatal(err)
						}
						lastStatus = !lastStatus

						time.Sleep(*flapInterval)
					}
				}(instanceID)
			}
		}

		log.Println("Services registered")
	}

	log.Println("Starting", *watchers, "watchers on", *serviceName, "...")

	for watcherID := 0; watcherID < *watchers; watcherID++ {
		go func() {
			lastIdx := uint64(0)

			for {
				_, meta, err := c.Health().Service(*serviceName, "", false, &consul.QueryOptions{
					WaitTime:   10 * time.Minute,
					WaitIndex:  lastIdx,
					AllowStale: true,
				})
				if err != nil {
					log.Fatal(err)
				}

				lastIdx = meta.LastIndex
				fmt.Print(".")
			}
		}()
	}

	log.Println("Watchers started.")

	close(startCh)

	<-make(chan bool)
}
