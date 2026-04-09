package game

import (
	"fmt"
	"math/rand"
)

const (
	JungleSize  = 15
	InnerSize   = 5
	CavernSizeX = 3
	CavernSizeY = 8
)

var QuestTargets = []string{"Capybara Den", "Ruins", "Waterfall", "Cavern"}

var InnerLocations = map[string]string{
	"1,1": "Village Elder",
	"3,1": "Shopkeeper",
	"1,3": "Frontiersman",
	"3,3": "Return",
	"2,1": "Basket Weaver",
	"2,3": "Guest Hut",
}

func NewGame() GameState {
	state := GameState{
		Player:      Player{X: 7, Y: 7, HP: 100, Energy: 100, Gold: 100},
		Mode:        ModeJungle,
		Discovered:  map[string]bool{},
		Locations:   map[string]Location{},
		Hazards:     map[string]string{},
		Camps:       map[string]Camp{},
		Inventory:   [3]*InventorySlot{},
		Quests:      []Quest{},
		BasketCap:   300,
		CavernTiles: map[string]string{},
		LastEvent:   "",
	}
	state.Locations["7,7"] = Location{Name: "Village", Enter: true}
	Randomize(&state)
	return state
}

func Randomize(s *GameState) {
	taken := map[string]bool{"7,7": true}

	for _, q := range QuestTargets {
		for {
			x := rand.Intn(JungleSize)
			y := rand.Intn(JungleSize)
			key := fmt.Sprintf("%d,%d", x, y)
			if !taken[key] {
				s.Locations[key] = Location{Name: q}
				taken[key] = true
				break
			}
		}
	}

	for _, h := range []string{"Beehive", "Jaguar", "Quicksand", "Anaconda"} {
		for {
			x := rand.Intn(JungleSize)
			y := rand.Intn(JungleSize)
			key := fmt.Sprintf("%d,%d", x, y)
			if !taken[key] {
				s.Hazards[key] = h
				taken[key] = true
				break
			}
		}
	}
}
