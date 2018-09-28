package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cakturk/gonetstat"
)

var (
	udp       = flag.Bool("udp", false, "display UDP sockets")
	tcp       = flag.Bool("tcp", false, "display TCP sockets")
	listening = flag.Bool("lis", false, "display only listening sockets")
	all       = flag.Bool("all", false, "display both listening and non-listening sockets")
)

// NetFlags represents the type of a flag
type NetFlags uint

// Different flag types
const (
	Listening NetFlags = iota + 1
	All
)

func main() {
	flag.Parse()

	var f NetFlags
	switch {
	case *all:
		f = All
	case *listening:
		f = Listening
	}

	if os.Geteuid() != 0 {
		fmt.Println("Not all processes could be identified, you would have to be root to see it all.")
	}
	fmt.Printf("Proto %-23s %-23s %-12s %-16s\n", "Local Addr", "Foreign Addr", "State", "PID/Program name")

	if *udp {
		tabs, err := gonetstat.UDPSocks()
		if err == nil {
			displaySockInfo("udp", f, tabs)
		}
	} else {
		*tcp = true
	}

	if *tcp {
		tabs, err := gonetstat.TCPSocks()
		if err == nil {
			displaySockInfo("tcp", f, tabs)
		}
	}
}

func displaySockInfo(proto string, f NetFlags, s []gonetstat.SockTabEntry) {
	for _, e := range s {
		switch f {
		case Listening:
			if e.State != 0x0a {
				continue
			}
		case All: // noop case
		default:
			if e.State == 0x0a {
				continue
			}
		}
		p := ""
		if e.Process != nil {
			p = e.Process.String()
		}

		fmt.Printf("%s   %-23s %-23s %-12s %-16s\n", proto, e.LocalAddr, e.RemoteAddr, e.State, p)
	}
}
