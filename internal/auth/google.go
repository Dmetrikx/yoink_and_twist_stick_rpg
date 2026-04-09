package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	sqlitedb "jungle-rpg/internal/repository/sqlite"
)

type GoogleAuth struct {
	config  *oauth2.Config
	queries *sqlitedb.Queries
}

func NewGoogleAuth(db *sql.DB) *GoogleAuth {
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/auth/google/callback"
	}
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email"},
		Endpoint:     google.Endpoint,
	}
	return &GoogleAuth{
		config:  cfg,
		queries: sqlitedb.New(db),
	}
}

func (g *GoogleAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := g.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (g *GoogleAuth) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, `{"error":"missing code"}`, http.StatusBadRequest)
		return
	}

	token, err := g.config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"token exchange failed: %s"}`, err), http.StatusInternalServerError)
		return
	}

	client := g.config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, `{"error":"failed to get user info"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		http.Error(w, `{"error":"failed to parse user info"}`, http.StatusInternalServerError)
		return
	}

	// Upsert user
	err = g.queries.UpsertUser(r.Context(), sqlitedb.UpsertUserParams{
		ID:    userInfo.ID,
		Email: userInfo.Email,
	})
	if err != nil {
		http.Error(w, `{"error":"failed to save user"}`, http.StatusInternalServerError)
		return
	}

	// Write session
	sess, _ := GetSession(r)
	sess.Values["user_id"] = userInfo.ID
	sess.Values["email"] = userInfo.Email
	if err := sess.Save(r, w); err != nil {
		http.Error(w, `{"error":"failed to save session"}`, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
