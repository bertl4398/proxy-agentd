package main

import (
  "log"
  "strings"
  "strconv"
  "os/exec"
)

func parsePktstatT(output []byte) (int, []string){
  output_str := string(output)
	lines := strings.Split(output_str, "\n")
	header := strings.Fields(lines[0])
	num_flows, _ := strconv.Atoi(header[0])
	// run_sec, _ := strconv.ParseFloat(header[1], 64)
	return num_flows, lines[1:len(lines)-1]
}

func Pktstat(iface string) (int, []string){
	out, err := exec.Command("pktstat", "-1nP", "-i", iface).Output()
	if err != nil {
		log.Fatal(err)
	}
	num_flows, flows := parsePktstatT(out)
	log.Printf("Number of flows: %d", num_flows)

  return num_flows, flows
}
