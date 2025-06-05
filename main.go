package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/gboncoffee/hearts/koro"
)

func main() {
	var k koro.KoroContext

	var dealer bool
	var peerAddress string
	var username string
	var peerPort int
	var localPort int
	var testMode bool

	flag.BoolVar(&dealer, "dealer", false, "dealer mode")
	flag.StringVar(&peerAddress, "pa", "localhost", "peer address")
	flag.StringVar(&username, "u", "", "username")
	flag.IntVar(&peerPort, "pp", koro.PORT, "peer port")
	flag.IntVar(&localPort, "lp", koro.PORT, "local port")
	flag.BoolVar(&testMode, "test", false, "test mode (autoplay)")

	flag.Parse()

	err := k.Init(peerAddress, peerPort, localPort, dealer)
	if err != nil {
		panic(err)
	}
	defer k.Fini()

	if dealer {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press enter to start...")
		reader.ReadString('\n')
	}

	peers := k.AssignNames(username, dealer)

	var game gameState
	game.players = peers
	game.points = make(map[koro.Address]int)
	game.dealer = dealer
	game.testMode = testMode

	i := 0
	for a := range peers {
		game.points[a] = 0
		game.order[i] = a
		i++
	}

	game.start(&k, dealer)
}
