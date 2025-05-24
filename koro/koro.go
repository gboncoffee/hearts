package koro

type Address uint8
type sequency uint8

type KoroContext struct {
	conn      connection
	names     map[Address]string
	outcoming chan Message
	incoming  chan Message
	addr      Address
}

func (k *KoroContext) Init(peerAddr string, peerPort int, localPort int) error {
	k.conn.init(256)

	err := k.conn.connectToPeer(peerAddr, peerPort)
	if err != nil {
		return err
	}
	err = k.conn.listen(localPort)
	if err != nil {
		return err
	}

	return nil
}

func (k *KoroContext) Fini() {
	k.conn.close()
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

	return parse(buffer)
}

func (k *KoroContext) Get() (Message, bool) {
	for {
		msg := k.get()
		if _, isYield := msg.(*tokenPassMessage); isYield {
			return nil, true
		}

		dest := msg.destination()
		k.send(msg)
		if dest == 0 || dest == k.addr {
			return msg, false
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
	if k.get() != msg {
		panic("network ack error")
	}
}

func (k *KoroContext) Yield() {
	k.send(&tokenPassMessage{})
}

func (k *KoroContext) AssignNames(username string, dealer bool) map[Address]string {
	k.addr = username2address(username)

	rightToSpeak := dealer

	for len(k.names) < 4 {
		if rightToSpeak {
			k.Send(&UsernameMessage{username: username}, 0)
			k.names[k.addr] = username
			k.Yield()
		}

		var msg Message
		msg, rightToSpeak = k.Get()
		if m, ok := msg.(*UsernameMessage); ok {
			k.names[username2address(m.username)] = m.username
		} else {
			panic("error assigning names")
		}
	}

	return nil
}
