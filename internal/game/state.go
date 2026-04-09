package game

type Mode string

const (
	ModeJungle Mode = "jungle"
	ModeInner  Mode = "inner"
	ModeCavern Mode = "cavern"
)

type Player struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	HP     int `json:"hp"`
	Energy int `json:"energy"`
	Gold   int `json:"gold"`
}

type InventorySlot struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Quest struct {
	Target string `json:"target"`
	Done   bool   `json:"done"`
}

type Camp struct {
	Uses int `json:"uses"`
}

type GameState struct {
	Player       Player              `json:"player"`
	Mode         Mode                `json:"mode"`
	Discovered   map[string]bool     `json:"discovered"`
	Locations    map[string]Location `json:"locations"`
	Hazards      map[string]string   `json:"hazards"`
	Camps        map[string]Camp     `json:"camps"`
	Inventory    [3]*InventorySlot   `json:"inventory"`
	Quests       []Quest             `json:"quests"`
	BasketCap    int                 `json:"basketCap"`
	CavernTiles  map[string]string   `json:"cavernTiles"`
	LastEvent    string              `json:"lastEvent"`
	CauseOfDeath string              `json:"causeOfDeath"`
	GameOver     bool                `json:"gameOver"`
	Won          bool                `json:"won"`
}

type Location struct {
	Name  string `json:"name"`
	Enter bool   `json:"enter,omitempty"`
}
