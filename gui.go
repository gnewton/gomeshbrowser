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
	log.Println("New session:", s.Id())
	s.SetTimeout(time.Second * 12)
	win := makeWin(h.db)

	winHandler := new(winHandler)
	win.AddEHandler(winHandler, gwu.ETypeWinLoad, gwu.ETypeStateChange, gwu.ETypeWinUnload, gwu.ETypeChange)

	s.AddWin(win)
	sessionWindowMap[s.Id()] = win
}
func (h sessHandler) Removed(s gwu.Session) {
	log.Println("Closing session:", s.Id())

	if win, ok := sessionWindowMap[s.Id()]; ok {
		removed := s.RemoveWin(win)
		log.Println("Removed window:", removed, win, "for session", s.Id())
	}

}

func makeWin(db *gorm.DB) gwu.Window {
	// Create and build a window
	win := gwu.NewWindow("main", "Test GUI Window")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HALeft)
	win.SetCellPadding(2)

	win.Add(gwu.NewSessMonitor())
	reset := gwu.NewLink("Reset", "/guitest/main")
	reset.SetTarget("_self")

	win.Add(reset)
	p := gwu.NewPanel()
	win.Add(p)
	l := gwu.NewLabel("PubMeSHicator")
	l.Style().SetColor(gwu.ClrGreen)
	p.Add(l)
	p.AddVSpace(5)

	//numChildren, topLevel, err := getTopLevel(db)
	_, topLevel, err := getLevel(0, db)
	if err != nil {
		log.Fatal(err)
	}
	for i, t := range topLevel {
		log.Println(i, t)
	}

	for _, t := range topLevel {
		numChildren := countChildren(t, db, true)

		meh := makeHandler(t, db)
		linePanel := gwu.NewHorizontalPanel()
		newExpander := gwu.NewExpander()
		newExpander.SetHeader(linePanel)
		makeExpanderContents(linePanel, t, numChildren)

		//e.SetHeader(gwu.NewLabel(findLabel(t) + " " + t.DescriptorName + " " + strconv.FormatInt(numChildren, 10) + " descendents"))
		newExpander.AddEHandler(meh, gwu.ETypeStateChange)
		p.Add(newExpander)
		p.AddVSpace(15)
	}
	return win
}
