package main

import (
	"log"
	"flag"
)

func logTraffic(iface string, c chan int) {
	num_flows, flows := Pktstat(iface)
	if num_flows > 0 {
		WriteRedisFlowBatch(flows)
	}
	c <- num_flows
}

func main() {
	var iface string
	var socket string
	flag.StringVar(&iface, "i", "lo", "capture interface")
	flag.StringVar(&socket, "s", "cmdsrv__0", "command socket")
	flag.Parse()

	InitLocalRedis()
	defer RedisConn.Close()

	log.Printf("Start capturing traffic at: %s", iface)
	flow_chan := make(chan int)
	go logTraffic(iface, flow_chan)

	log.Printf("Start command server at: %s", socket)
	go StartUnixDomainSocketServer(socket)
	defer StopUnixDomainSocketServer(socket)

	for {
		select {
		case <- flow_chan:
			go logTraffic(iface, flow_chan)
		}
	}
}
