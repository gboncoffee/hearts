package main

import (
	"flag"
	"time"

	"github.com/gboncoffee/hearts/koro"
)

func main() {
	var k koro.KoroContext

	var dealer bool
	var peerAddress string
	var username string
	var peerPort int
	var localPort int

	flag.BoolVar(&dealer, "dealer", false, "dealer mode")
	flag.StringVar(&peerAddress, "pa", "localhost", "peer address")
	flag.StringVar(&username, "u", "", "username")
	flag.IntVar(&peerPort, "pp", koro.PORT, "peer port")
	flag.IntVar(&localPort, "lp", koro.PORT, "local port")

	flag.Parse()

	err := k.Init("localhost", peerPort, localPort)
	if err != nil {
		panic(err)
	}

	k.AssignNames(username, dealer)

	time.Sleep(time.Second * 2)

	k.Fini()
}
