package main

import (
  "net"
  "sync"
  "strings"
  "bufio"
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
  log.WithFields(log.Fields{
    "port": CONN_PORT,
    "protocol": CONN_TYPE,
  }).Info("server listening")

  for {
    conn, err := l.Accept()
    if err != nil {
      log.Error(err)
      continue
    }
    defer conn.Close()
    go handleRequest(conn)
  }
}

func handleRequest(c net.Conn) {
  scanner := bufio.NewScanner(c)
  for scanner.Scan() {
    data := scanner.Text()
    log.WithFields(log.Fields{
      "data": data,
    }).Debug("data received")
    go executeCmd(data)
  }
  if err := scanner.Err(); err != nil {
    log.Error(err)
  }
}

func executeCmd(data string) {
  var j map[string]interface{}
  json.Unmarshal([]byte(data), &j)
  cmd := j["message"].(string)

  switch {
  case strings.HasPrefix(cmd, "BLK"):
    fields := strings.Fields(cmd)
    if len(fields) == 4 {
      proto := fields[1]
      ip := fields[2]
      port := fields[3]
      log.WithFields(log.Fields{
        "ip": ip,
        "port": port,
        "action": "blocK",
      }).Info(cmd)
      BlockIpAtPort(proto, ip, port)
    }
  case strings.HasPrefix(cmd, "RDR"):
    fields := strings.Fields(cmd)
    if len(fields) == 2 {
      ip := fields[1]
      log.WithFields(log.Fields{
        "ip": ip,
        "action": "redirect",
      }).Info(cmd)
      RedirectIp(ip)
    }
  default:
    log.Info(cmd)
  }
}
