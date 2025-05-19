package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var conn connection
	var serverMode bool
	var port int
	var peerPort int

	flag.BoolVar(&serverMode, "serverMode", false, "set's server mode")
	flag.IntVar(&port, "port", PORT, "override default port")
	flag.IntVar(&peerPort, "peerPort", PORT, "override default peer port")
	flag.Parse()

	conn.listen(port)

	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("enter peer address: ")
	reader.Scan()
	addr := reader.Text()

	conn.connectToPeer(addr, peerPort)

	if serverMode {
		for {
			msg := conn.read()
			fmt.Printf("received message: %s\n", string(msg))
		}
	} else {
		for {
			conn.send([]byte("Hello, World!"))
		}
	}
}
