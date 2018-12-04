package main

import (
	"flag"
	"log"
	"os"
	"time"

	consul "github.com/hashicorp/consul/api"
)

func main() {
	consulAddr := flag.String("consul", "127.0.0.1:8500", "Consul address")
	serviceName := flag.String("service", "srv", "Service to watch")
	registerInstances := flag.Int("register", 0, "Register N -service instances")
	deregister := flag.Bool("deregister", false, "Deregister all instances of -service")
	flapInterval := flag.Duration("flap-interval", 0, "If -register is given, flap each instance between critical and passing state on given interval")
	token := flag.String("token", "", "ACL token")
	watchers := flag.Int("watchers", 1, "Number of concurrnet watchers on service")
	flag.Parse()

	if *token == "" {
		*token = os.Getenv("ACL_TOKEN")
	}

	c, err := consul.NewClient(&consul.Config{
		Address: *consulAddr,
		Token:   *token,
	})
	if err != nil {
		log.Fatal(err)
	}

	stats := make(chan Stat)

	if *registerInstances > 0 {
		err := RegisterServices(c, *serviceName, *registerInstances, *flapInterval, stats)
		if err != nil {
			log.Fatal(err)
		}
	} else if *deregister {
		err := DeregisterServices(c, *serviceName)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	go RunQueries(QueryAgent(c, *serviceName, 10*time.Minute, true), *watchers, stats)

	DisplayStats(stats)
}
