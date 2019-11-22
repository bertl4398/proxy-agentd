package main

import (
	"flag"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	RedirectEndpoint string
}

var Config Configuration

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	var logfile string
	var database string
	var configfile string

	flag.StringVar(&configfile, "c", "conf.json", "path to configuration file")
	flag.StringVar(&database, "d", "test", "influxDB database name")
	flag.StringVar(&logfile, "f", "/var/log/ulog/gprint.log", "path to log file")
	flag.Parse()

	Config = Configuration{}
	if err := gonfig.GetConf(configfile, &Config); err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	log.Info("Start logging traffic to influxDB")
	wg.Add(1)
	go StartLogTraffic(logfile, database, wg)
	wg.Add(1)
	go StartTcpSocketServer(wg)

	wg.Wait()
}
