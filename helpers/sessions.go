package helpers

var dbSessions map[string]string = make(map[string]string)

func GetSession(sid string) string {
	sessionValue, ok := dbSessions[sid]
	if ok {
		return sessionValue
	}

	return ""
}

func SetSession(sid string, value string) {
	dbSessions[sid] = value
}

func UnsetSession(sid string) {
	delete(dbSessions, sid)
}
