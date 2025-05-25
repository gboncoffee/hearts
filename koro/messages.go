package koro

import "fmt"

type messageType uint8

const (
	nack = iota
	username_message
	your_cards_message
	play_message
	token_pass_message
)

type Message interface {
	origin() Address
	destination() Address
	setOrigin(Address)
	setDestination(Address)
}

type commonMessage struct {
	orig Address
	dest Address
}

func (c *commonMessage) destination() Address {
	return c.dest
}

func (c *commonMessage) origin() Address {
	return c.orig
}

func (c *commonMessage) setOrigin(orig Address) {
	c.orig = orig
}

func (c *commonMessage) setDestination(dest Address) {
	c.dest = dest
}

type UsernameMessage struct {
	commonMessage
	username string
}

type YourCardsMessage struct {
	commonMessage
	Cards [13]byte
}

type PlayMessage struct {
	commonMessage
	card byte
}

type tokenPassMessage struct {
	commonMessage
}

func parse(buffer []byte) Message {
	orig := buffer[0]
	dest := buffer[1]
	tipe := buffer[2]

	cm := commonMessage{Address(orig), Address(dest)}

	switch tipe {
	case username_message:
		size := buffer[3]
		return &UsernameMessage{
			commonMessage: cm,
			username:      string(buffer[4 : 4+size]),
		}
	case your_cards_message:
		msg := YourCardsMessage{commonMessage: cm}
		for i, c := range buffer[3 : 3+13] {
			msg.Cards[i] = byte(c)
		}
		return &msg
	case play_message:
		return &PlayMessage{
			commonMessage: cm,
			card:          buffer[3],
		}
	case token_pass_message:
		return &tokenPassMessage{
			commonMessage: cm,
		}
	default:
		panic(fmt.Errorf("unknown message type %v", tipe))
	}
}

func serialize(msg Message) []byte {
	var tipe messageType
	var data []byte
	switch m := msg.(type) {
	case *UsernameMessage:
		tipe = username_message
		data = []byte{byte(len(m.username))}
		data = append(data, []byte(m.username)...)
	case *YourCardsMessage:
		tipe = your_cards_message
		data = make([]byte, len(m.Cards))
		for i, c := range m.Cards {
			data[i] = byte(c)
		}
	case *PlayMessage:
		tipe = play_message
		data = []byte{byte(m.card)}
	case *tokenPassMessage:
		tipe = token_pass_message
	}

	bin := []byte{
		byte(msg.origin()),
		byte(msg.destination()),
		byte(tipe),
	}
	return append(bin, data...)
}
