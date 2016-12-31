package main

import (
	"log"

	"github.com/icza/gowut/gwu"
)

func main() {
	//log.SetOutput(ioutil.Discard)

	db, err := dbOpen("mesh2016_sqlite3.db")
	if err != nil {
		log.Fatal(err)
	}

	handlerMap = make(map[int64]*meshExpanderHandler)
	sessionWindowMap = make(map[string]gwu.Window)

	server := gwu.NewServer("guitest", "localhost:8081")
	server.AddSessCreatorName("main", "Login Window")
	server.AddSHandler(sessHandler{db: db})
	server.Start("") // Also opens windows list in browser
}
