package main

import (
  "sync"
  "os/exec"

  log "github.com/sirupsen/logrus"
)

var rules = struct{
  sync.RWMutex
  m map[string]int
}{m: make(map[string]int)}


func ListIptables() {
  out, err := exec.Command("iptables", "-S").Output()
  if err != nil {
      log.Error(err)
  }
  log.WithFields(log.Fields{
    "out": out,
  }).Info("iptables -S")
}

func BlockIpAtPort(proto string, ip string, port string) {
  rules.Lock()
  if _, ok := rules.m["BLK"+ip]; ok {
    rules.m["BLK"+ip]++
  } else {
    args := []string{"-A", "INPUT", "-s", ip, "-p", proto, "--dport", port, "-j", "DROP"}
    _, err := exec.Command("iptables", args...).Output()
    if err != nil {
        log.Error(err)
    } else {
      rules.m["BLK"+ip]++
    }
  }

  n := rules.m["BLK"+ip]
  log.WithFields(log.Fields{
    "ip": ip,
    "count": n,
  }).Info("block")

  rules.Unlock()
}

func RedirectIp(ip string) {
  rules.Lock()
  if _, ok := rules.m["RDR"+ip]; ok {
    rules.m["RDR"+ip]++
  } else {
    args := []string{"-t", "nat", "-A", "PREROUTING", "-s", ip, "-p", "tcp",
      "--dport", "80", "-j", "DNAT", "--to-destination", "192.168.122.197:80"}
    _, err := exec.Command("iptables", args...).Output()
    if err != nil {
        log.Error(err)
    } else {
      args = []string{"-t", "nat", "-A", "POSTROUTING", "-j", "MASQUERADE"}
      _, err = exec.Command("iptables", args...).Output()
      if err != nil {
        log.Error(err)
      }
      rules.m["RDR"+ip]++
    }
  }

  n := rules.m["RDR"+ip]
  log.WithFields(log.Fields{
  "rule": ip,
  "count": n,
  }).Info("redirect")

  rules.Unlock()
}
