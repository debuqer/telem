package helpers

import "github.com/google/uuid"

func GetCsrfToken(sid string) string {
	token := uuid.Must(uuid.NewRandom()).String()
	SetSession(sid, "csrf", token)

	return token
}

func MatchCsrf(sid string, token string) bool {
	if csrfToken, hasItem := GetSession(sid, "csrf"); hasItem {
		return token == csrfToken
	}

	return false
}
