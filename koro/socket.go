package koro

import (
	"net"
)

const PORT = 42069

type connection struct {
	peerConn *net.UDPConn
	conn     *net.UDPConn
	bufSize  int
}

func (c *connection) init(bufSize int) {
	c.bufSize = bufSize
}

func (c *connection) connectToPeer(addr string, port int) error {
	peerAddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(addr),
	}

	conn, err := net.DialUDP("udp", nil, &peerAddr)
	if err != nil {
		return err
	}

	c.peerConn = conn
	return nil
}

func (c *connection) listen(port int) error {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *connection) close() {
	c.peerConn.Close()
	c.conn.Close()
}

func (c *connection) send(msg []byte) {
	c.peerConn.Write(msg)
}

func (c *connection) read() ([]byte, error) {
	buffer := make([]byte, c.bufSize)
	n, _, err := c.conn.ReadFromUDP(buffer)
	return buffer[:n], err
}
