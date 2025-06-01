package main

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/gboncoffee/hearts/koro"
)

type suit uint8
type card = uint8

const (
	SPADES = iota
	HEARTS
	CLUBS
	DIAMONDS
)

const (
	SPADES_A = iota
	SPADES_2
	SPADES_3
	SPADES_4
	SPADES_5
	SPADES_6
	SPADES_7
	SPADES_8
	SPADES_9
	SPADES_10
	SPADES_J
	SPADES_Q
	SPADES_K
	HEARTS_A
	HEARTS_2
	HEARTS_3
	HEARTS_4
	HEARTS_5
	HEARTS_6
	HEARTS_7
	HEARTS_8
	HEARTS_9
	HEARTS_10
	HEARTS_J
	HEARTS_Q
	HEARTS_K
	CLUBS_A
	CLUBS_2
	CLUBS_3
	CLUBS_4
	CLUBS_5
	CLUBS_6
	CLUBS_7
	CLUBS_8
	CLUBS_9
	CLUBS_10
	CLUBS_J
	CLUBS_Q
	CLUBS_K
	DIAMONDS_A
	DIAMONDS_2
	DIAMONDS_3
	DIAMONDS_4
	DIAMONDS_5
	DIAMONDS_6
	DIAMONDS_7
	DIAMONDS_8
	DIAMONDS_9
	DIAMONDS_10
	DIAMONDS_J
	DIAMONDS_Q
	DIAMONDS_K
)

func getSuit(c card) suit {
	if c < HEARTS_A {
		return SPADES
	}
	if c < CLUBS_A {
		return HEARTS
	}
	if c < DIAMONDS_A {
		return CLUBS
	}
	return DIAMONDS
}

type gameState struct {
	players map[koro.Address]string
	points  map[koro.Address]int
	dealer  bool
}

func deal(k *koro.KoroContext, players [4]koro.Address) *[13]card {
	cards := [52]card{}
	for i := range 52 {
		cards[i] = card(i)
	}

	rand.Shuffle(len(cards), func(i, j int) {
		tmp := cards[i]
		cards[i] = cards[j]
		cards[j] = tmp
	})

	we := k.Address()
	ourCards := new([13]card)
	for i, player := range players {
		msg := koro.YourCardsMessage{}
		playerCards := &msg.Cards
		if player == we {
			playerCards = ourCards
		}

		j := 0
		for k := i * 13; k < (i+1)*13; k++ {
			playerCards[j] = cards[i]
			j++
		}

		if player != we {
			k.Send(&msg, player)
		}
	}

	return ourCards
}

func waitDeal(k *koro.KoroContext) *[13]card {
	for {
		msg, _ := k.Get()
		return (*[13]card)(msg.(*koro.YourCardsMessage).Cards)
	}
}

type rankEntry struct {
	player string
	points int
}

func (g *gameState) finish() *[4]rankEntry {
	rank := new([4]rankEntry)
	i := 0
	someoneExploded := false
	for player, points := range g.points {
		entry := rankEntry{
			player: g.players[player],
			points: points,
		}
		rank[i] = entry
		i++
		if points >= 100 {
			someoneExploded = true
		}
	}
	if !someoneExploded {
		return nil
	}

	slices.SortStableFunc(rank[:], func(a rankEntry, b rankEntry) int {
		if a.points < b.points {
			return -1
		}
		if a.points > b.points {
			return 1
		}
		return 0
	})

	return rank
}

func (g *gameState) start(k *koro.KoroContext, dealer bool) {
	rank := g.mainloop(k, dealer)
	fmt.Println("End!")
	fmt.Println("Ranking:")
	for i, e := range rank {
		fmt.Printf("%v. %-20s: %v", i, e.player, e.points)
	}
}

func (g *gameState) mainloop(k *koro.KoroContext, dealer bool) *[4]rankEntry {
	for {
		g.round(k, dealer)
		rank := g.finish()
		if rank != nil {
			return rank
		}
	}
}

type roundState struct {
	hand   []card
	points map[koro.Address]int
	broke  bool
}

func (g *gameState) round(k *koro.KoroContext, dealer bool) {
	var players [4]koro.Address
	i := 0
	for a := range g.players {
		players[i] = a
		i++
	}

	var cards *[13]card
	if dealer {
		cards = deal(k, players)
	} else {
		cards = waitDeal(k)
	}

	hand := cards[:]
	for {
		// TODO
	}
}
