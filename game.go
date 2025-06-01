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

const UNDEFINED_SUIT = 0xff

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

func card2string(c card) string {
	switch c {
	case SPADES_A:
		return "🂡"
	case SPADES_2:
		return "🂢"
	case SPADES_3:
		return "🂣"
	case SPADES_4:
		return "🂤"
	case SPADES_5:
		return "🂥"
	case SPADES_6:
		return "🂦"
	case SPADES_7:
		return "🂧"
	case SPADES_8:
		return "🂨"
	case SPADES_9:
		return "🂩"
	case SPADES_10:
		return "🂪"
	case SPADES_J:
		return "🂫"
	case SPADES_Q:
		return "🂭"
	case SPADES_K:
		return "🂮"
	case HEARTS_A:
		return "🂱"
	case HEARTS_2:
		return "🂲"
	case HEARTS_3:
		return "🂳"
	case HEARTS_4:
		return "🂴"
	case HEARTS_5:
		return "🂵"
	case HEARTS_6:
		return "🂶"
	case HEARTS_7:
		return "🂷"
	case HEARTS_8:
		return "🂸"
	case HEARTS_9:
		return "🂹"
	case HEARTS_10:
		return "🂺"
	case HEARTS_J:
		return "🂻"
	case HEARTS_Q:
		return "🂽"
	case HEARTS_K:
		return "🂾"
	case CLUBS_A:
		return "🃑"
	case CLUBS_2:
		return "🃒"
	case CLUBS_3:
		return "🃓"
	case CLUBS_4:
		return "🃔"
	case CLUBS_5:
		return "🃕"
	case CLUBS_6:
		return "🃖"
	case CLUBS_7:
		return "🃗"
	case CLUBS_8:
		return "🃘"
	case CLUBS_9:
		return "🃙"
	case CLUBS_10:
		return "🃚"
	case CLUBS_J:
		return "🃛"
	case CLUBS_Q:
		return "🃝"
	case CLUBS_K:
		return "🃞"
	case DIAMONDS_A:
		return "🃁"
	case DIAMONDS_2:
		return "🃂"
	case DIAMONDS_3:
		return "🃃"
	case DIAMONDS_4:
		return "🃄"
	case DIAMONDS_5:
		return "🃅"
	case DIAMONDS_6:
		return "🃆"
	case DIAMONDS_7:
		return "🃇"
	case DIAMONDS_8:
		return "🃈"
	case DIAMONDS_9:
		return "🃉"
	case DIAMONDS_10:
		return "🃊"
	case DIAMONDS_J:
		return "🃋"
	case DIAMONDS_Q:
		return "🃍"
	case DIAMONDS_K:
		return "🃎"
	default:
		panic("unknown card")
	}
}

type gameState struct {
	order   [4]koro.Address
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
		return (&msg.(*koro.YourCardsMessage).Cards)
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

type trickState struct {
	table map[koro.Address]card
	suit  suit
}

func (g *gameState) round(k *koro.KoroContext, dealer bool) {
	var cards *[13]card
	if dealer {
		cards = deal(k, g.order)
	} else {
		cards = waitDeal(k)
	}

	fmt.Println("Round starting!")
	fmt.Println("Your hand:")
	for _, c := range cards {
		fmt.Printf("%v ", card2string(c))
	}
	fmt.Println()

	round := roundState{
		hand:   cards[:],
		points: make(map[koro.Address]int),
		broke:  false,
	}

	for a := range g.points {
		round.points[a] = 0
	}

	rightToSpeak := dealer
	start := slices.Index(cards[:], CLUBS_2) != -1

	// Actually play the round.
	for len(round.hand) > 0 {
		start, rightToSpeak, round.broke = g.trick(
			k,
			&round,
			start,
			rightToSpeak,
			round.broke,
		)
	}

	for a := range g.points {
		g.points[a] += round.points[a]
	}
}

func (g *gameState) trick(
	k *koro.KoroContext,
	round *roundState,
	start bool,
	rightToSpeak bool,
	broke bool,
) (bool, bool, bool) {
	trick := trickState{
		table: make(map[koro.Address]card),
		suit:  suit(UNDEFINED_SUIT),
	}

	if start {
		for !rightToSpeak {
			_, rightToSpeak = k.Get()
		}
		shouldBreak := g.play(k, round, &trick, broke, true)
		if !broke {
			broke = shouldBreak
		}
	} else if rightToSpeak {
		k.Yield()
	}

	ourTime := false
	var msg koro.Message
	for len(trick.table) != 4 {
		if ourTime {
			shouldBreak := g.play(k, round, &trick, broke, false)
			if !broke {
				broke = shouldBreak
			}
			k.Yield()
			continue
		}

		msg, ourTime = k.Get()
		if ourTime && len(trick.table) == 0 {
			k.Yield()
			ourTime = false
			continue
		}

		card := card(msg.(*koro.PlayMessage).Card)
		trick.table[msg.Origin()] = card
		if trick.suit == UNDEFINED_SUIT {
			suit := getSuit(card)
			if !broke {
				broke = suit == HEARTS
			}

			trick.suit = suit
		}
	}

	winner := koro.Address(0)
	winnerCard := card(0)
	points := 0
	for a := range trick.table {
		suit := getSuit(trick.table[a])
		if suit == HEARTS {
			points++
		} else if trick.table[a] == SPADES_Q {
			points += 13
		}
		if suit == trick.suit {
			if trick.table[a] >= winnerCard {
				winnerCard = trick.table[a]
				winner = a
			}
		}
	}

	round.points[winner] += points

	return winner == k.Address(), ourTime, broke
}

func (g *gameState) play(
	k *koro.KoroContext,
	round *roundState,
	trick *trickState,
	broke bool,
	first bool,
) bool {
	allowed := getAllowedCards(trick.suit, round.hand, broke, first)

	fmt.Println("Your time!")
	fmt.Println("Your entire hand: ")
	for _, c := range round.hand {
		fmt.Printf("%v ", card2string(c))
	}
	fmt.Println()
	fmt.Println("You're allowed to play:")
	for i, c := range allowed {
		fmt.Printf("%v - %v \n", i, c)
	}
	fmt.Print("Choose a card to play: ")
	var n int
	fmt.Scanf("%d", &n)

	for n > 0 && n <= len(allowed) {
		fmt.Printf("\nChoose a number between 1 and %v: ", len(allowed))
		fmt.Scanf("%d", &n)
	}

	card := allowed[n-1]
	idx := slices.Index(round.hand, card)
	round.hand = slices.Delete(round.hand, idx, idx+1)
	k.Send(&koro.PlayMessage{Card: card}, 0)
	k.Yield()
	trick.table[k.Address()] = card
	suit := getSuit(card)
	if first {
		trick.suit = suit
	}

	return suit == HEARTS
}

func getAllowedCards(suit suit, hand []card, broke bool, first bool) (allowed []card) {
	hasSuit := false
	if suit != UNDEFINED_SUIT {
		for _, c := range hand {
			if getSuit(c) == suit {
				hasSuit = true
				break
			}
		}
	}

	if !hasSuit {
		if first && !broke {
			hasNotHearts := false
			for _, c := range hand {
				if getSuit(c) != HEARTS {
					hasNotHearts = true
					break
				}
			}
			if !hasNotHearts {
				return hand
			} else {
				// We remove all hearts from the hand, keeping every other suit.
				for _, c := range hand {
					if getSuit(c) != HEARTS {
						allowed = append(allowed, c)
					}
				}
				return
			}
		}
		return hand
	}

	// We remove everything that's not from the suit.
	for _, c := range hand {
		if getSuit(c) == suit {
			allowed = append(allowed, c)
		}
	}

	return
}
