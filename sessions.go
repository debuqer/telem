package main

var dbSessions map[string]string = make(map[string]string)

func getSession(sid string) string {
	sessionValue, ok := dbSessions[sid]
	if ok {
		return sessionValue
	}

	return ""
}

func setSession(sid string, value string) {
	dbSessions[sid] = value
}

func unsetSession(sid string) {
	delete(dbSessions, sid)
}
