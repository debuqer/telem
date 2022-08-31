package helpers

import "net/http"

func GetCookieValue(c *http.Cookie, err error) string {
	if err == http.ErrNoCookie {
		return ""
	}

	return c.Value
}

func SetCookie(w http.ResponseWriter, name string, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:  name,
		Value: value,
	})
}
