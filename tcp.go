package main

import (
  "net"
  "sync"
  "strings"
  "io/ioutil"
  "encoding/json"

  log "github.com/sirupsen/logrus"
)

const (
  CONN_HOST = "localhost"
  CONN_PORT = "9988"
  CONN_TYPE = "tcp"
)

func StartTcpSocketServer(wg *sync.WaitGroup) {
  defer wg.Done()
  l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
  if err != nil {
    log.Fatal(err)
  }
  defer l.Close()
  log.Info("Start listening on " + CONN_HOST + ":" + CONN_PORT)

  for {
    conn, err := l.Accept()
    if err != nil {
      log.Fatal(err)
    }
    go handleRequest(conn)
  }
}

func handleRequest(c net.Conn) {
  for {
    buf, err := ioutil.ReadAll(c)
    if err != nil {
      log.Error(err)
    }
    var result map[string]interface{}
    json.Unmarshal([]byte(buf), &result)

    cmd := result["message"].(string)
    executeCmd(cmd)
  }
}

func executeCmd(cmd string) {
  switch {
  case strings.HasPrefix(cmd, "BLK"):
    fields := strings.Fields(cmd)
    if len(fields) == 4 {
      // proto := fields[1]
      ip := fields[2]
      port := fields[3]
      log.Printf("Block IP %s from port %s", ip, port)
    }
  case strings.HasPrefix(cmd, "RDR"):
    fields := strings.Fields(cmd)
    if len(fields) == 2 {
      ip := fields[1]
      log.Printf("Redirect IP: %s", ip)
      RedirectIp(ip)
    }
  default:
    log.Info("Received command ", cmd)
  }
}
