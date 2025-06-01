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
	defer k.Fini()

	if dealer {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press enter to start...")
		reader.ReadString('\n')
	}

	peers := k.AssignNames(username, dealer)
	fmt.Println("Game starting!")
	fmt.Print("Players connected:")
	for _, p := range peers {
		fmt.Printf(" %v", p)
	}
	fmt.Println()

	var game gameState
	game.players = peers
	game.points = make(map[koro.Address]int)
	game.dealer = dealer

	for a := range peers {
		game.points[a] = 0
	}

	game.start(&k, dealer)
}
