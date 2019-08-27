package main

import (
  "os"
  "time"
  "sync"
  "strings"
  "strconv"

  "github.com/hpcloud/tail"
  log "github.com/sirupsen/logrus"
  _ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
  client "github.com/influxdata/influxdb1-client/v2"
)

func WriteGprint(s string, db string) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	var t time.Time
	var srcport, dstport, srcip, dstip, srcmac string
	var csum int

	key_value := strings.Split(s, ",")
	for _, x := range key_value {
		kv := strings.Split(x, "=")
		switch key := kv[0]; key {
		case "timestamp":
			// 2019/08/21-07:32:40
			t, _ = time.Parse("2006/01/02-15:04:05", kv[1])
		case "tcp.sport":
			srcport = kv[1]
		case "tcp.dport":
			dstport = kv[1]
		case "ip.saddr":
			srcip = kv[1]
		case "ip.daddr":
			dstip = kv[1]
		case "mac.saddr.str":
			srcmac = kv[1]
		case "tcp.csum":
			csum, _ = strconv.Atoi(kv[1])
		default:
		}
	}
  if csum <= 0 {
    log.Warn("unknown entry")
    return
  }

	tags := map[string]string{
		"ip.saddr": srcip,
	  "ip.daddr": dstip,
		"mac.saddr": srcmac,
		"tcp.sport": srcport,
		"tcp.dport": dstport,
	}
	fields := map[string]interface{}{
		"tcp.csum": csum,
	}
	pt, err := client.NewPoint("flow", tags, fields, t)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}

  // log.Info(pt.String())
}

func StartLogTraffic(logfile string, database string, wg *sync.WaitGroup) {
  defer wg.Done()
  t, _ := tail.TailFile(logfile, tail.Config{
    Logger: log.StandardLogger(),
    Location: &tail.SeekInfo{0, os.SEEK_END},
    Follow: true,
    ReOpen: true})
  for line := range t.Lines {
    s := line.Text
    WriteGprint(s, database)
  }
}
