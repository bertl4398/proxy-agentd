package main

import (
  "os/exec"

  log "github.com/sirupsen/logrus"
)

func ListIptables() {
  out, err := exec.Command("iptables", "-S").Output()
  if err != nil {
      log.Fatal(err)
  }
  log.Printf("%s", out)
}

func BlockIpAtPort(proto string, ip string, port string) {
  log.Printf("iptables -A INPUT -s %s -p %s --dport %s -j DROP",
    ip, proto, port)
  out, err := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-p", proto,
    "--dport", port, "-j", "DROP").Output()
  if err != nil {
      log.Error(err)
  }
  log.Printf("%s", out)
}

func RedirectIp(ip string) {
  log.Printf("iptables -t nat -A PREROUTING -s %s -p tcp --dport 80 -j DNAT" +
    "--to-destination 192.168.122.197:80", ip)
  _, err := exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-s",
    ip, "-p", "tcp", "--dport", "80", "-j", "DNAT", "--to-destination",
    "192.168.122.197:80").Output()
  if err != nil {
      log.Error(err)
  }
  log.Printf("iptables -t nat -A POSTROUTING -j MASQUERADE")
  _, err = exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING",
    "-j", "MASQUERADE").Output()
  if err != nil {
    log.Error(err)
  }
}
