package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"jungle-rpg/internal/auth"
	sqlitedb "jungle-rpg/internal/repository/sqlite"
)

type SaveHandler struct {
	queries *sqlitedb.Queries
}

func NewSaveHandler(db *sql.DB) *SaveHandler {
	return &SaveHandler{queries: sqlitedb.New(db)}
}

func (h *SaveHandler) ListSaves(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r)
	saves, err := h.queries.ListSaves(r.Context(), userID)
	if err != nil {
		jsonError(w, "failed to list saves", http.StatusInternalServerError)
		return
	}
	if saves == nil {
		saves = []sqlitedb.ListSavesRow{}
	}
	jsonResponse(w, saves)
}

func (h *SaveHandler) CreateSave(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r)
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	state, err := loadGameFromSession(r)
	if err != nil {
		jsonError(w, "no active game to save", http.StatusBadRequest)
		return
	}
	stateJSON, _ := json.Marshal(state)

	save, err := h.queries.CreateSave(r.Context(), sqlitedb.CreateSaveParams{
		UserID: userID,
		Name:   req.Name,
		State:  stateJSON,
	})
	if err != nil {
		jsonError(w, "failed to create save", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonResponse(w, save)
}

func (h *SaveHandler) UpdateSave(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r)
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, "invalid save id", http.StatusBadRequest)
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	state, err := loadGameFromSession(r)
	if err != nil {
		jsonError(w, "no active game to save", http.StatusBadRequest)
		return
	}
	stateJSON, _ := json.Marshal(state)

	err = h.queries.UpdateSave(r.Context(), sqlitedb.UpdateSaveParams{
		Name:   req.Name,
		State:  stateJSON,
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		jsonError(w, "failed to update save", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *SaveHandler) LoadSave(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r)
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, "invalid save id", http.StatusBadRequest)
		return
	}

	save, err := h.queries.GetSave(r.Context(), sqlitedb.GetSaveParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		jsonError(w, "save not found", http.StatusNotFound)
		return
	}

	// Load state into session
	sess, _ := auth.GetSession(r)
	sess.Values["game_state"] = string(save.State)
	if err := sess.Save(r, w); err != nil {
		jsonError(w, "failed to load save into session", http.StatusInternalServerError)
		return
	}

	var state json.RawMessage = save.State
	jsonResponse(w, state)
}

func (h *SaveHandler) DeleteSave(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r)
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		jsonError(w, "invalid save id", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteSave(r.Context(), sqlitedb.DeleteSaveParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		jsonError(w, "failed to delete save", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
