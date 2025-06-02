package koro

import "fmt"

type Address uint8
type sequency uint8

type KoroContext struct {
	conn         connection
	names        map[Address]string
	outcoming    chan Message
	incoming     chan Message
	addr         Address
	rightToSpeak bool
}

func (k *KoroContext) Init(peerAddr string, peerPort int, localPort int, rts bool) error {
	k.conn.init(256)

	err := k.conn.connectToPeer(peerAddr, peerPort)
	if err != nil {
		return err
	}
	err = k.conn.listen(localPort)
	if err != nil {
		return err
	}

	k.rightToSpeak = rts

	return nil
}

func (k *KoroContext) Fini() {
	k.conn.close()
}

func (k *KoroContext) Address() Address {
	return k.addr
}

func username2address(name string) Address {
	hash := uint8(0)
	for _, b := range []uint8(name) {
		hash += b
	}
	return Address(hash)
}

func (k *KoroContext) get() Message {
	buffer, err := k.conn.read()
	if err != nil {
		panic(err)
	}

	pmsg := parse(buffer)
	return pmsg
}

func (k *KoroContext) Get() Message {
	for {
		msg := k.get()
		if _, isYield := msg.(*tokenPassMessage); isYield {
			k.rightToSpeak = true
			return nil
		}

		dest := msg.destination()
		k.send(msg)
		if dest == 0 || dest == k.addr {
			return msg
		}
	}
}

func (k *KoroContext) send(msg Message) {
	bin := serialize(msg)
	k.conn.send(bin)
}

func (k *KoroContext) Send(msg Message, dest Address) {
	msg.setOrigin(k.addr)
	msg.setDestination(dest)
	k.send(msg)
	k.get()
}

func (k *KoroContext) Yield() {
	if !k.rightToSpeak {
		panic("Called yield without the right to speak.")
	}
	k.send(&tokenPassMessage{})
	k.rightToSpeak = false
}

func (k *KoroContext) RightToSpeak() bool {
	return k.rightToSpeak
}

func (k *KoroContext) AssignNames(username string, rightToSpeak bool) map[Address]string {
	k.addr = username2address(username)
	k.names = make(map[Address]string)
	k.rightToSpeak = rightToSpeak

	for {
		if len(k.names) == 4 {
			// If we did start with the rts, we should get it back.
			if rightToSpeak {
				for !k.RightToSpeak() {
					k.Get()
				}
			}
			return k.names
		}
		if k.RightToSpeak() {
			k.Send(&UsernameMessage{username: username}, 0)
			k.names[k.addr] = username
			k.Yield()
			continue
		}

		var msg Message
		msg = k.Get()
		if msg != nil {
			if m, ok := msg.(*UsernameMessage); ok {
				k.names[username2address(m.username)] = m.username
			} else {
				panic(fmt.Errorf("got %T when assign names", msg))
			}
		}
	}
}
