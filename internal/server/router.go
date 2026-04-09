package server

import (
	"database/sql"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"jungle-rpg/internal/api"
	"jungle-rpg/internal/auth"
)

func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	googleAuth := auth.NewGoogleAuth(db)
	gameHandler := api.NewGameHandler()
	saveHandler := api.NewSaveHandler(db)

	// Auth routes
	r.Get("/auth/google", googleAuth.HandleLogin)
	r.Get("/auth/google/callback", googleAuth.HandleCallback)

	// API routes (require auth)
	r.Route("/api", func(r chi.Router) {
		r.Use(auth.RequireAuth)

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			sess, _ := auth.GetSession(r)
			email, _ := sess.Values["email"].(string)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"email": email})
		})

		// Game routes
		r.Route("/game", func(r chi.Router) {
			r.Post("/new", gameHandler.NewGame)
			r.Get("/state", gameHandler.GetState)
			r.Post("/move", gameHandler.Move)
			r.Post("/action", gameHandler.Action)
			r.Post("/item", gameHandler.UseItem)
		})

		// Save routes
		r.Route("/saves", func(r chi.Router) {
			r.Get("/", saveHandler.ListSaves)
			r.Post("/", saveHandler.CreateSave)
			r.Put("/{id}", saveHandler.UpdateSave)
			r.Get("/{id}/load", saveHandler.LoadSave)
			r.Delete("/{id}", saveHandler.DeleteSave)
		})
	})

	// Serve React static files
	webDist := findWebDist()
	if webDist != "" {
		fileServer(r, webDist)
	}

	return r
}

func findWebDist() string {
	// Check relative to working directory
	candidates := []string{
		"web/dist",
		"./web/dist",
	}
	// Check relative to executable
	ex, err := os.Executable()
	if err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(ex), "web", "dist"))
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}

func fileServer(r chi.Router, root string) {
	fsys := os.DirFS(root)
	fileServer := http.FileServer(http.FS(fsys))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Check if file exists
		if _, err := fs.Stat(fsys, path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for non-API routes
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
