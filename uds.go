package main

import (
  "os"
  "net"
  "log"
  "strings"
)

func cmdSrv(c net.Conn) {
  for {
    buf := make([]byte, 1024)
    nr, err := c.Read(buf)
    if err != nil {
      return // EOF
    }
    data := buf[0:nr]
    log.Printf("Server got: %s", data)
    // _, err = c.Write(data)
    // if err != nil {
    //     log.Fatal(err)
    // }

    cmd := string(data)
    if strings.HasPrefix(cmd, "BLK") {
      fields := strings.Fields(cmd)
      if len(fields) == 4 {
        proto := fields[1]
        ip := fields[2]
        port := fields[3]
        log.Printf("Block IP %s from port %s", ip, port)
        BlockIpAtPort(proto, ip, port)
      }
    }
    if strings.HasPrefix(cmd, "RDR") {
      fields := strings.Fields(cmd)
      ip := fields[1]
      log.Printf("Redirect IP: %s", ip)
    }
  }
}

func StartUnixDomainSocketServer(socket string) {
  StopUnixDomainSocketServer(socket)
  l, err := net.Listen("unix", socket)
  if err != nil {
    log.Fatal(err)
  }
  os.Chown(socket, 1000, 1000)

  log.Printf("listening on socket %s", socket)
  for {
    fd, err := l.Accept()
    if err != nil {
      log.Fatal(err)
    }
    go cmdSrv(fd)
  }
}

func StopUnixDomainSocketServer(socket string) {
  if _, err := os.Stat(socket); err == nil {
    err := os.Remove(socket)
    if err != nil {
      log.Fatal(err)
    }
  }
}
