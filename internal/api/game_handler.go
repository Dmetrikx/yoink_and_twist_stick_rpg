package api

import (
	"encoding/json"
	"net/http"

	"jungle-rpg/internal/auth"
	"jungle-rpg/internal/game"
)

type GameHandler struct{}

func NewGameHandler() *GameHandler {
	return &GameHandler{}
}

func (h *GameHandler) NewGame(w http.ResponseWriter, r *http.Request) {
	state := game.NewGame()
	if err := saveGameToSession(r, w, state); err != nil {
		jsonError(w, "failed to save game state", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, state)
}

func (h *GameHandler) GetState(w http.ResponseWriter, r *http.Request) {
	state, err := loadGameFromSession(r)
	if err != nil {
		jsonError(w, "no active game", http.StatusNotFound)
		return
	}
	jsonResponse(w, state)
}

func (h *GameHandler) Move(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DX int `json:"dx"`
		DY int `json:"dy"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	state, err := loadGameFromSession(r)
	if err != nil {
		jsonError(w, "no active game", http.StatusNotFound)
		return
	}
	state, _ = game.MovePlayer(state, req.DX, req.DY)
	if err := saveGameToSession(r, w, state); err != nil {
		jsonError(w, "failed to save game state", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, state)
}

func (h *GameHandler) Action(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	state, err := loadGameFromSession(r)
	if err != nil {
		jsonError(w, "no active game", http.StatusNotFound)
		return
	}
	state, _ = game.DoAction(state, req.Action)
	if err := saveGameToSession(r, w, state); err != nil {
		jsonError(w, "failed to save game state", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, state)
}

func (h *GameHandler) UseItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Slot int `json:"slot"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	state, err := loadGameFromSession(r)
	if err != nil {
		jsonError(w, "no active game", http.StatusNotFound)
		return
	}
	state, _ = game.UseItem(state, req.Slot)
	if err := saveGameToSession(r, w, state); err != nil {
		jsonError(w, "failed to save game state", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, state)
}

// Session helpers

func saveGameToSession(r *http.Request, w http.ResponseWriter, state game.GameState) error {
	sess, err := auth.GetSession(r)
	if err != nil {
		return err
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	sess.Values["game_state"] = string(data)
	return sess.Save(r, w)
}

func loadGameFromSession(r *http.Request) (game.GameState, error) {
	sess, err := auth.GetSession(r)
	if err != nil {
		return game.GameState{}, err
	}
	raw, ok := sess.Values["game_state"].(string)
	if !ok || raw == "" {
		return game.GameState{}, http.ErrNoCookie
	}
	var state game.GameState
	err = json.Unmarshal([]byte(raw), &state)
	return state, err
}

// JSON helpers

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
