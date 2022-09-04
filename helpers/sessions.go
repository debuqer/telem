package helpers

var dbSessions map[string]string = make(map[string]string)

func GetSession(sid string) (string, bool) {
	var val string
	var hasItem bool = false

	conn, err := GetConn()
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	stmt, err := conn.Prepare("SELECT value FROM sessions WHERE sid = ?")
	if err != nil {
		panic(err)
	}
	row, err := stmt.Query(sid)

	for row.Next() {
		row.Scan(&val)
		hasItem = true
	}

	return val, hasItem
}

func SetSession(sid string, value string) {
	conn, err := GetConn()
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	_, hasItem := GetSession(sid)
	if hasItem {
		stmt, err := conn.Prepare("UPDATE sessions SET value = ? WHERE sid = ?")
		if err != nil {
			panic(err)
		}

		stmt.Exec(value, sid)
	} else {
		stmt, err := conn.Prepare("INSERT INTO sessions ( sid, value, created_at, expires_at) VALUES( ?, ?, NOW(), NOW())")
		if err != nil {
			panic(err)
		}
		stmt.Exec(sid, value)
	}
}

func UnsetSession(sid string) {
	conn, err := GetConn()
	defer conn.Close()

	stmt, err := conn.Prepare("DELETE FROM sessions WHERE sid = ?")
	if err != nil {
		panic(err)
	}

	stmt.Exec(sid)
}
