package game

import (
	"fmt"
	"math/rand"
)

// MovePlayer attempts to move the player by (dx, dy).
func MovePlayer(s GameState, dx, dy int) (GameState, string) {
	if s.GameOver || s.Won {
		return s, ""
	}

	switch s.Mode {
	case ModeJungle:
		return moveJungle(s, dx, dy)
	case ModeInner:
		return moveInner(s, dx, dy)
	case ModeCavern:
		return moveCavern(s, dx, dy)
	}
	return s, ""
}

func moveJungle(s GameState, dx, dy int) (GameState, string) {
	nx := s.Player.X + dx
	ny := s.Player.Y + dy
	if nx < 0 || nx >= JungleSize || ny < 0 || ny >= JungleSize {
		return s, "You can't go that way."
	}

	s.Player.X = nx
	s.Player.Y = ny
	key := fmt.Sprintf("%d,%d", nx, ny)
	s.Discovered[key] = true

	// Energy cost
	if s.Player.Energy > 0 {
		s.Player.Energy -= 8
		if s.Player.Energy < 0 {
			s.Player.Energy = 0
		}
	} else {
		s.Player.HP -= 8
	}

	// Check death from energy drain
	if s.Player.HP <= 0 {
		s.GameOver = true
		if s.CauseOfDeath == "" {
			s.CauseOfDeath = "Exhaustion"
		}
		s.LastEvent = fmt.Sprintf("You died of %s!", s.CauseOfDeath)
		return s, s.LastEvent
	}

	// Hazards
	if h, ok := s.Hazards[key]; ok {
		event := resolveHazard(&s, h)
		if s.Player.HP <= 0 {
			s.GameOver = true
			s.LastEvent = fmt.Sprintf("You died! %s", event)
			return s, s.LastEvent
		}
		s.LastEvent = event
		return s, event
	}

	// Locations
	if loc, ok := s.Locations[key]; ok {
		event := fmt.Sprintf("You found: %s", loc.Name)
		// Check quest completion
		for i, q := range s.Quests {
			if !q.Done && q.Target == loc.Name {
				s.Quests[i].Done = true
				reward := 550
				if loc.Name == "Waterfall" {
					reward = 700
				}
				s.Player.Gold += reward
				s.BasketCap += 300
				event += fmt.Sprintf(" — Quest complete! +%d Gold", reward)
			}
		}
		s.LastEvent = event
		return s, event
	}

	// Camp
	if camp, ok := s.Camps[key]; ok && camp.Uses > 0 {
		s.LastEvent = fmt.Sprintf("You found your camp (%d uses left).", camp.Uses)
		return s, s.LastEvent
	}

	// Random events
	event := randomEvent(&s)
	if s.Player.HP <= 0 {
		s.GameOver = true
		s.LastEvent = fmt.Sprintf("You died! %s", event)
		return s, s.LastEvent
	}
	if event != "" {
		s.LastEvent = event
		return s, event
	}

	s.LastEvent = "You push through the jungle..."
	return s, s.LastEvent
}

func moveInner(s GameState, dx, dy int) (GameState, string) {
	nx := s.Player.X + dx
	ny := s.Player.Y + dy
	if nx < 0 || nx >= InnerSize || ny < 0 || ny >= InnerSize {
		return s, "You can't go that way."
	}
	s.Player.X = nx
	s.Player.Y = ny
	key := fmt.Sprintf("%d,%d", nx, ny)

	if name, ok := InnerLocations[key]; ok {
		s.LastEvent = fmt.Sprintf("You approach the %s.", name)
		return s, s.LastEvent
	}
	s.LastEvent = "You walk through the village."
	return s, s.LastEvent
}

func moveCavern(s GameState, dx, dy int) (GameState, string) {
	nx := s.Player.X + dx
	ny := s.Player.Y + dy
	if nx < 0 || nx >= CavernSizeX || ny < 0 || ny >= CavernSizeY {
		return s, "You can't go that way."
	}
	s.Player.X = nx
	s.Player.Y = ny
	key := fmt.Sprintf("%d,%d", nx, ny)

	// Check exit
	if key == "1,0" {
		s.Mode = ModeJungle
		// Restore jungle position from locations
		for k, loc := range s.Locations {
			if loc.Name == "Cavern" {
				var cx, cy int
				fmt.Sscanf(k, "%d,%d", &cx, &cy)
				s.Player.X = cx
				s.Player.Y = cy
				break
			}
		}
		s.CavernTiles = map[string]string{}
		s.LastEvent = "You exit the cavern."
		return s, s.LastEvent
	}

	// Cavern tiles
	if tile, ok := s.CavernTiles[key]; ok {
		switch tile {
		case "Deep Hole":
			s.Player.HP -= 50
			s.Player.Energy -= 15
			if s.Player.Energy < 0 {
				s.Player.Energy = 0
			}
			event := "You fall into a deep hole! -50 HP -15 Energy"
			if s.Player.HP <= 0 {
				s.GameOver = true
				s.CauseOfDeath = "Deep Hole"
				s.LastEvent = fmt.Sprintf("You died! %s", event)
				return s, s.LastEvent
			}
			s.LastEvent = event
			return s, event
		case "Hidden Treasure":
			s.Player.Gold += 400
			delete(s.CavernTiles, key)
			s.LastEvent = "Hidden Treasure! +400 Gold"
			return s, s.LastEvent
		case "Hidden Spring":
			s.Player.HP += 75
			s.Player.Energy += 100
			delete(s.CavernTiles, key)
			s.LastEvent = "Hidden Spring! +75 HP +100 Energy"
			return s, s.LastEvent
		}
	}

	s.LastEvent = "You explore the dark cavern..."
	return s, s.LastEvent
}

func resolveHazard(s *GameState, hazard string) string {
	switch hazard {
	case "Beehive":
		s.Player.HP -= 25
		s.CauseOfDeath = "Bees"
		return "Beehive! A swarm attacks you! -25 HP"
	case "Jaguar":
		s.Player.HP -= 65
		s.CauseOfDeath = "Jaguar"
		return "Jaguar! It pounces on you! -65 HP"
	case "Anaconda":
		s.Player.HP = 5
		s.CauseOfDeath = "Anaconda"
		return "Anaconda! It constricts you! HP set to 5"
	case "Quicksand":
		s.Player.Energy -= 45
		if s.Player.Energy < 0 {
			s.Player.Energy = 0
		}
		return "Quicksand! You struggle free. -45 Energy"
	}
	return ""
}

func randomEvent(s *GameState) string {
	roll := rand.Intn(100)
	switch {
	case roll < 5:
		s.Player.HP -= 25
		s.CauseOfDeath = "Slip and Fall"
		return "Slip and Fall! -25 HP"
	case roll < 10:
		s.Player.Gold += 30
		return "This Looks Valuable! +30 Gold"
	case roll < 25:
		s.Player.HP -= 10
		s.CauseOfDeath = "Twisted Ankle"
		return "Twisted Ankle! -10 HP"
	case roll < 40:
		s.Player.Gold += 10
		return "What's This? +10 Gold"
	default:
		return ""
	}
}

// DoAction performs a named action at the player's current tile.
func DoAction(s GameState, action string) (GameState, string) {
	if s.GameOver || s.Won {
		return s, ""
	}

	switch action {
	case "enterVillage":
		return enterVillage(s)
	case "exitVillage":
		return exitVillage(s)
	case "enterCavern":
		return enterCavern(s)
	case "exitCavern":
		return exitCavern(s)
	case "rest":
		return rest(s)
	case "buildCamp":
		return buildCamp(s)
	case "addQuest":
		return addQuest(s)
	case "buyPotion":
		return buyItem(s, "Health Potion", 50)
	case "buyStim":
		return buyItem(s, "Stimulant", 25)
	case "basketWeave":
		return basketWeave(s)
	case "buyVillage":
		return buyVillage(s)
	}
	return s, "Unknown action."
}

func enterVillage(s GameState) (GameState, string) {
	key := fmt.Sprintf("%d,%d", s.Player.X, s.Player.Y)
	loc, ok := s.Locations[key]
	if !ok || loc.Name != "Village" {
		return s, "No village here."
	}
	s.Mode = ModeInner
	s.Player.X = 2
	s.Player.Y = 2
	s.LastEvent = "You enter the village."
	return s, s.LastEvent
}

func exitVillage(s GameState) (GameState, string) {
	if s.Mode != ModeInner {
		return s, "You're not in the village."
	}
	s.Mode = ModeJungle
	// Return to village tile on the jungle map
	for k, loc := range s.Locations {
		if loc.Name == "Village" {
			var vx, vy int
			fmt.Sscanf(k, "%d,%d", &vx, &vy)
			s.Player.X = vx
			s.Player.Y = vy
			break
		}
	}
	s.LastEvent = "You leave the village."
	return s, s.LastEvent
}

func enterCavern(s GameState) (GameState, string) {
	key := fmt.Sprintf("%d,%d", s.Player.X, s.Player.Y)
	loc, ok := s.Locations[key]
	if !ok || loc.Name != "Cavern" {
		return s, "No cavern here."
	}
	s.Mode = ModeCavern
	s.Player.X = 1
	s.Player.Y = 7
	s.CavernTiles = generateCavern()
	s.LastEvent = "You descend into the cavern..."
	return s, s.LastEvent
}

func generateCavern() map[string]string {
	tiles := map[string]string{}
	taken := map[string]bool{
		"1,7": true, // entrance
		"1,0": true, // exit
	}
	items := []string{"Deep Hole", "Hidden Treasure", "Hidden Spring"}
	for _, item := range items {
		for {
			x := rand.Intn(CavernSizeX)
			y := rand.Intn(CavernSizeY)
			key := fmt.Sprintf("%d,%d", x, y)
			if !taken[key] {
				tiles[key] = item
				taken[key] = true
				break
			}
		}
	}
	return tiles
}

func exitCavern(s GameState) (GameState, string) {
	if s.Mode != ModeCavern {
		return s, "You're not in a cavern."
	}
	s.Mode = ModeJungle
	for k, loc := range s.Locations {
		if loc.Name == "Cavern" {
			var cx, cy int
			fmt.Sscanf(k, "%d,%d", &cx, &cy)
			s.Player.X = cx
			s.Player.Y = cy
			break
		}
	}
	s.CavernTiles = map[string]string{}
	s.LastEvent = "You exit the cavern."
	return s, s.LastEvent
}

func rest(s GameState) (GameState, string) {
	key := fmt.Sprintf("%d,%d", s.Player.X, s.Player.Y)

	// Guest Hut in village
	if s.Mode == ModeInner {
		if name, ok := InnerLocations[key]; ok && name == "Guest Hut" {
			s.Player.Energy += 30
			s.Player.HP += 20
			s.LastEvent = "You rest at the Guest Hut. +30 Energy +20 HP"
			return s, s.LastEvent
		}
		return s, "Nothing to do here."
	}

	// Camp in jungle
	if camp, ok := s.Camps[key]; ok && camp.Uses > 0 {
		camp.Uses--
		s.Camps[key] = camp
		s.Player.Energy += 30
		s.Player.HP += 20
		s.LastEvent = fmt.Sprintf("You rest at camp. +30 Energy +20 HP (%d uses left)", camp.Uses)
		return s, s.LastEvent
	}
	return s, "Nothing to rest at here."
}

func buildCamp(s GameState) (GameState, string) {
	if s.Mode != ModeJungle {
		return s, "You can only build camps in the jungle."
	}
	if s.Player.Gold < 50 {
		return s, "Not enough gold! Need 50."
	}
	key := fmt.Sprintf("%d,%d", s.Player.X, s.Player.Y)
	if _, ok := s.Locations[key]; ok {
		return s, "Can't build a camp at a location."
	}
	if _, ok := s.Camps[key]; ok {
		return s, "There's already a camp here."
	}
	s.Player.Gold -= 50
	s.Camps[key] = Camp{Uses: 2}
	s.LastEvent = "Camp built! (-50 Gold, 2 uses)"
	return s, s.LastEvent
}

func addQuest(s GameState) (GameState, string) {
	if s.Mode != ModeInner {
		return s, "Talk to the Village Elder inside the village."
	}
	key := fmt.Sprintf("%d,%d", s.Player.X, s.Player.Y)
	name, ok := InnerLocations[key]
	if !ok || name != "Village Elder" {
		return s, "The Village Elder is not here."
	}

	// Find unassigned quest targets
	assigned := map[string]bool{}
	for _, q := range s.Quests {
		assigned[q.Target] = true
	}
	for _, target := range QuestTargets {
		if !assigned[target] {
			s.Quests = append(s.Quests, Quest{Target: target, Done: false})
			s.LastEvent = fmt.Sprintf("New quest: Find the %s!", target)
			return s, s.LastEvent
		}
	}
	s.LastEvent = "No more quests available."
	return s, s.LastEvent
}

func buyItem(s GameState, itemName string, cost int) (GameState, string) {
	if s.Mode != ModeInner {
		return s, "You need to be in the village."
	}
	key := fmt.Sprintf("%d,%d", s.Player.X, s.Player.Y)
	name, ok := InnerLocations[key]
	if !ok || name != "Shopkeeper" {
		return s, "The Shopkeeper is not here."
	}
	if s.Player.Gold < cost {
		return s, fmt.Sprintf("Not enough gold! Need %d.", cost)
	}

	// Find existing slot or empty slot
	for i := range s.Inventory {
		if s.Inventory[i] != nil && s.Inventory[i].Name == itemName {
			if s.Inventory[i].Count >= 5 {
				return s, fmt.Sprintf("You already have the max (5) %ss.", itemName)
			}
			s.Player.Gold -= cost
			s.Inventory[i].Count++
			s.LastEvent = fmt.Sprintf("Bought %s! (-%d Gold)", itemName, cost)
			return s, s.LastEvent
		}
	}
	for i := range s.Inventory {
		if s.Inventory[i] == nil {
			s.Player.Gold -= cost
			s.Inventory[i] = &InventorySlot{Name: itemName, Count: 1}
			s.LastEvent = fmt.Sprintf("Bought %s! (-%d Gold)", itemName, cost)
			return s, s.LastEvent
		}
	}
	return s, "Inventory full!"
}

func basketWeave(s GameState) (GameState, string) {
	if s.Mode != ModeInner {
		return s, "You need to be in the village."
	}
	key := fmt.Sprintf("%d,%d", s.Player.X, s.Player.Y)
	name, ok := InnerLocations[key]
	if !ok || name != "Basket Weaver" {
		return s, "The Basket Weaver is not here."
	}
	if s.Player.Gold >= s.BasketCap {
		return s, "Your basket is full!"
	}
	s.Player.Gold += 6
	if s.Player.Gold > s.BasketCap {
		s.Player.Gold = s.BasketCap
	}
	s.LastEvent = "The Basket Weaver pays you. +6 Gold"
	return s, s.LastEvent
}

func buyVillage(s GameState) (GameState, string) {
	if s.Mode != ModeInner {
		return s, "You need to be in the village."
	}
	key := fmt.Sprintf("%d,%d", s.Player.X, s.Player.Y)
	name, ok := InnerLocations[key]
	if !ok || name != "Frontiersman" {
		return s, "The Frontiersman is not here."
	}
	if s.Player.Gold < 1500 {
		return s, "Not enough gold! Need 1500."
	}
	s.Player.Gold -= 1500
	s.Won = true
	s.LastEvent = "You purchased the village! You win!"
	return s, s.LastEvent
}

// UseItem uses the inventory item at slot index i (0-2).
func UseItem(s GameState, i int) (GameState, string) {
	if s.GameOver || s.Won {
		return s, ""
	}
	if i < 0 || i >= 3 || s.Inventory[i] == nil {
		return s, "No item in that slot."
	}

	slot := s.Inventory[i]
	switch slot.Name {
	case "Health Potion":
		s.Player.HP += 20
		s.LastEvent = "Used Health Potion! +20 HP"
	case "Stimulant":
		s.Player.Energy += 10
		s.LastEvent = "Used Stimulant! +10 Energy"
	default:
		return s, "Can't use that."
	}

	slot.Count--
	if slot.Count <= 0 {
		s.Inventory[i] = nil
	} else {
		s.Inventory[i] = slot
	}
	return s, s.LastEvent
}
