package middleware

import (
	"net/http"

	"github.com/SOG-web/goinit/gin/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// NewSessionMiddleware creates and returns a session middleware
func NewSessionMiddleware(cfg config.Config) gin.HandlerFunc {
	store := cookie.NewStore([]byte(cfg.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   cfg.SessionDomain,
		MaxAge:   cfg.SessionMaxAge,
		Secure:   cfg.SessionSecure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // or http.SameSiteStrictMode, http.SameSiteNoneMode
	})
	return sessions.Sessions(cfg.SessionName, store)
}

// GetSession returns the session from the context
func GetSession(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}

// SetSessionValue sets a value in the session
func SetSessionValue(c *gin.Context, key string, value interface{}) {
	session := GetSession(c)
	session.Set(key, value)
	session.Save()
}

// GetSessionValue gets a value from the session
func GetSessionValue(c *gin.Context, key string) interface{} {
	session := GetSession(c)
	return session.Get(key)
}

// DeleteSessionValue deletes a value from the session
func DeleteSessionValue(c *gin.Context, key string) {
	session := GetSession(c)
	session.Delete(key)
	session.Save()
}

// ClearSession clears all session data
func ClearSession(c *gin.Context) {
	session := GetSession(c)
	session.Clear()
	session.Save()
}