package main

import (
	"log"
	"time"

	"github.com/icza/gowut/gwu"
	"github.com/jinzhu/gorm"
)

type sessHandler struct {
	db *gorm.DB
}

type winHandler struct {
}

var sessionWindowMap map[string]gwu.Window

func (h *winHandler) HandleEvent(ev gwu.Event) {
	log.Println("New win event", ev)
}

func (h sessHandler) Created(s gwu.Session) {
	if hrr, ok := s.(gwu.HasRequestResponse); ok {
		req := hrr.Request()
		log.Println("Client addr:", req.RemoteAddr)
	}

	s.SetTimeout(time.Minute * (SESSION_LENGTH + 2))
	win := makeWin(h.db, s)

	winHandler := new(winHandler)
	win.AddEHandler(winHandler, gwu.ETypeWinLoad, gwu.ETypeStateChange, gwu.ETypeWinUnload, gwu.ETypeChange)

	s.AddWin(win)
	sessionWindowMap[s.Id()] = win
}
func (h sessHandler) Removed(s gwu.Session) {
	log.Println("Closing session:", s.Id())

	if win, ok := sessionWindowMap[s.Id()]; ok {
		removed := s.RemoveWin(win)
		delete(sessionWindowMap, s.Id())
		log.Println("Removed window:", removed, win, "for session", s.Id())
	}

}

func makeWin(db *gorm.DB, s gwu.Session) gwu.Window {
	// Create and build a window
	win := gwu.NewWindow("main", "Test GUI Window")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HALeft)
	win.SetCellPadding(2)

	//win.Add(gwu.NewSessMonitor())
	reset := gwu.NewLink(RESET_SESSION, "/guitest/main")
	reset.SetTarget("_self")

	p := gwu.NewPanel()
	win.Add(p)
	p2 := gwu.NewPanel()

	topPanel := gwu.NewHorizontalPanel()
	p.Add(topPanel)

	l := gwu.NewLabel(APP_TITLE)
	l.Style().SetColor(gwu.ClrGreen)
	topPanel.Add(l)
	topPanel.AddHSpace(500)

	timm := gwu.NewTimer(time.Minute * SESSION_LENGTH)
	win.Add(timm)

	timm.AddEHandlerFunc(func(e gwu.Event) {
		log.Println("TIMER: removing session")
		win.Remove(timm)
		win.Remove(p)
		win.Remove(p2)
		s.SetTimeout(0 * time.Minute)
		expiredSession := gwu.NewLink("Session expired", "/guitest/main")
		expiredSession.SetTarget("_self")
		win.Add(expiredSession)
		e.MarkDirty(win, p, p2, timm)
	}, gwu.ETypeStateChange)

	topPanel.Add(reset)

	_, topLevel, err := getLevel(0, db)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range topLevel {
		numChildren := countChildren(t, db, true)

		meh := makeHandler(t, db)
		linePanel := gwu.NewHorizontalPanel()
		newExpander := gwu.NewExpander()
		newExpander.SetHeader(linePanel)
		makeExpanderContents(linePanel, t, numChildren)

		newExpander.AddEHandler(meh, gwu.ETypeStateChange)
		p2.Add(newExpander)
		p2.AddVSpace(15)
	}
	p2.AddVSpace(75)
	l = gwu.NewLink("Glen Newton", "https://github.com/gnewton")
	l.Style().SetColor(gwu.ClrGreen).Set("text-decoration", "none")
	p2.Add(l)
	p2.AddVSpace(15)
	l = gwu.NewLink("gomeshbrowser@github", "https://github.com/gnewton/gomeshbrowser")
	l.Style().SetColor(gwu.ClrGreen).Set("text-decoration", "none")

	p2.Add(l)
	win.Add(p2)
	return win
}
