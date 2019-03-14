package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/criteo/consul-bench/stats"
	consul "github.com/hashicorp/consul/api"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

type step struct {
	Value    int
	Duration time.Duration
}

func main() {
	consulAddr := flag.String("consul", "127.0.0.1:8500", "Consul address")
	useRPC := flag.Bool("rpc", false, "Use RPC server calls instead of agent HTTP")
	rpcAddr := flag.String("rpc-addr", "127.0.0.1:8300", "When using rpc, the consul rpc addr")
	dc := flag.String("dc", "dc1", "When using rpc, the consul datacenter")
	serviceName := flag.String("service", "srv", "Service to watch")
	registerInstances := flag.Int("register", 0, "Register N -service instances")
	deregister := flag.Bool("deregister", false, "Deregister all instances of -service")
	flapInterval := flag.Duration("flap-interval", 0, "If -register is given, flap each instance between critical and passing state on given interval")
	wait := flag.Duration("query-wait", 10*time.Minute, "Bloquing queries max wait time")
	stale := flag.Bool("query-stale", false, "Run stale blocking queries")
	token := flag.String("token", "", "ACL token")
	watchers := flag.String("watchers", "1", "Number of concurrnet watchers on service")
	monitor := flag.Int("monitor", 0, "Consul PID")
	runtime := flag.Duration("time", 0, "Time to run the benchmark")
	latepc := flag.Float64("late-ratio", 0, "Ratio of late callers")
	makePlot := flag.Bool("plot", false, "Draw a QPS plot")
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

	statsC := make(chan stats.Stat)

	if *registerInstances > 0 {
		err := RegisterServices(c, *serviceName, *registerInstances, *flapInterval, statsC)
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

	done := make(chan struct{})
	var wg sync.WaitGroup

	if *monitor > 0 {
		wg.Add(1)
		go func() {
			Monitor(int32(*monitor), statsC, done)
			wg.Done()
		}()
	}

	if *runtime > 0 {
		go func() {
			time.Sleep(*runtime)
			close(done)
		}()
	}

	var qf queryFn
	if !*useRPC {
		qf = QueryAgent(c, *serviceName, *wait, *stale)
	} else {
		qf = QueryServer(*rpcAddr, *dc, *serviceName, *wait, *stale)
	}

	sr := stats.New()
	go sr.Run(statsC)

	steps, err := paseSteps(*watchers)
	if err != nil {
		log.Fatal(err)
	}

	type stepStat struct {
		Watchers int
		QPS      float64
	}

	stepsStats := []stepStat{}

	wg.Add(1)
	go func() {
		for _, step := range steps {
			wdone := make(chan struct{})
			timer := time.NewTimer(step.Duration)
			go func() {
				select {
				case <-done:
					close(wdone)
				case <-timer.C:
					close(wdone)
				}
			}()
			sr.Reset()
			RunQueries(qf, step.Value, *latepc, statsC, wdone)
			avgs := sr.AVGs()
			qps := 0.0
			for _, v := range avgs {
				if v.Label == "QPS" {
					qps = v.Value
				}
			}
			fmt.Println(sr.AVGs())
			stepsStats = append(stepsStats, stepStat{
				Watchers: step.Value,
				QPS:      qps,
			})
		}
		close(done)
		wg.Done()
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		<-signals
		close(done)
	}()

	<-done

	if *makePlot {
		p, err := plot.New()
		if err != nil {
			panic(err)
		}

		p.Title.Text = "QPS / Watchers"
		p.X.Label.Text = "Watchers"
		p.Y.Label.Text = "QPS"

		points := plotter.XYs{}
		for _, step := range stepsStats {
			points = append(points, plotter.XY{
				X: float64(step.Watchers),
				Y: step.QPS,
			})
		}

		err = plotutil.AddLinePoints(p,
			"A", points,
		)
		if err != nil {
			panic(err)
		}

		// Save the plot to a PNG file.
		if err := p.Save(800, 600, "points.png"); err != nil {
			panic(err)
		}
	}

	wg.Wait()
	os.Exit(0)
}

func paseSteps(raw string) ([]step, error) {
	steps := []step{}
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		i := strings.IndexByte(p, ':')
		if i == -1 {
			return nil, fmt.Errorf("invalid step %s", p)
		}

		value, err := strconv.Atoi(p[:i])
		if err != nil {
			return nil, fmt.Errorf("invalid step %s", p)
		}

		dur, err := time.ParseDuration(p[i+1:])
		if err != nil {
			return nil, fmt.Errorf("invalid step %s", p)
		}

		steps = append(steps, step{value, dur})
	}

	return steps, nil
}
