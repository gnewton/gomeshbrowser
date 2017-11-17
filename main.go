package main

import (
	"log"

	"github.com/icza/gowut/gwu"
)

func main() {
	//log.SetOutput(ioutil.Discard)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := dbOpen(DB_FILE)
	if err != nil {
		log.Fatal(err)
	}

	handlerMap = make(map[int64]*meshExpanderHandler)
	sessionWindowMap = make(map[string]gwu.Window)

	server := gwu.NewServer("guitest", "localhost:8081")
	server.AddSessCreatorName("main", APP_TITLE)
	server.AddSHandler(sessHandler{db: db})
	server.Start("") // Also opens windows list in browser
}
