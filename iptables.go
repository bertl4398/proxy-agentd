package main

import (
  "log"
  "os/exec"
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
      log.Fatal(err)
  }
  log.Printf("%s", out)
}
