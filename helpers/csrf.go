package helpers

import (
	"net/http"

	"github.com/google/uuid"
)

func GetCsrfToken(w http.ResponseWriter, r *http.Request) string {
	cookieToken := uuid.Must(uuid.NewRandom()).String()

	SetCookie(w, "csrf_token", cookieToken)

	token := uuid.Must(uuid.NewRandom()).String()
	SetSession(cookieToken, "csrf", token)

	return token
}

func MatchCsrf(r *http.Request, token string) bool {
	coockieToken := GetCookieValue(r.Cookie("csrf_token"))
	if csrfToken, hasItem := GetSession(coockieToken, "csrf"); hasItem {
		return token == csrfToken
	}

	return false
}
