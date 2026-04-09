package auth

import (
	"encoding/gob"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const sessionDir = "/tmp/jungle-sessions"

const sessionName = "jungle-session"

var store *sessions.FilesystemStore

func init() {
	gob.Register(map[string]interface{}{})
}

func InitStore() {
	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		panic("failed to create session directory: " + err.Error())
	}
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		key = "dev-session-key-change-in-prod!!"
	}
	store = sessions.NewFilesystemStore(sessionDir, []byte(key))
	store.MaxLength(1 << 20) // 1MB max session size for game state
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func GetStore() *sessions.FilesystemStore {
	return store
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, sessionName)
}

func GetUserID(r *http.Request) (string, error) {
	sess, err := GetSession(r)
	if err != nil {
		return "", err
	}
	uid, ok := sess.Values["user_id"].(string)
	if !ok || uid == "" {
		return "", errors.New("not authenticated")
	}
	return uid, nil
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := GetUserID(r)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
