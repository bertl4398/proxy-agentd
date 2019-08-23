package main

import (
	"flag"
	"sync"

	log "github.com/sirupsen/logrus"
)

func main() {
	var logfile string
	var database string

	flag.StringVar(&database, "d", "test", "influxDB database name")
	flag.StringVar(&logfile, "f", "/var/log/ulog/gprint.log", "path to log file")
	flag.Parse()

	wg := new(sync.WaitGroup)
	log.Info("Start logging traffic to influxDB")
	wg.Add(1)
	go StartLogTraffic(logfile, database, wg)
	wg.Add(1)
	go StartTcpSocketServer(wg)

	wg.Wait()
}
