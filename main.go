package main

import (
	"io/ioutil"
	"log"
	"strconv"

	"github.com/icza/gowut/gwu"
	"github.com/jinzhu/gorm"
)

//These should be shared across sessions
// Cache in hash
type meshExpanderHandler struct {
	meshTree *MeshTree
	db       *gorm.DB
}

type linkClickHandler struct {
}

var handlerMap map[int64]*meshExpanderHandler

func makeHandler(meshTree *MeshTree, db *gorm.DB) *meshExpanderHandler {
	if val, ok := handlerMap[meshTree.ID]; ok {
		log.Println("Cache hit")
		return val
	}
	log.Println("Cache miss")
	newHandler := &meshExpanderHandler{meshTree: meshTree, db: db}
	handlerMap[meshTree.ID] = newHandler

	return newHandler
}

func (h *linkClickHandler) HandleEvent(ev gwu.Event) {
	log.Println("***********************")
	log.Println(ev)
	log.Printf("FFFFFF-- event %+v\n", ev)
	log.Printf("FFFFFF-- event src %+v\n", ev.Src())
	log.Printf("FFFFFF-- parent %+v\n", ev.Parent())
	log.Printf("FFFFFF-- parent src %+v\n", ev.Parent().Src())
}

func (h *meshExpanderHandler) HandleEvent(ev gwu.Event) {
	if hrr, ok := ev.(gwu.HasRequestResponse); ok {
		req := hrr.Request()
		log.Println("Client addr:", req.RemoteAddr)
		log.Println(ev)

	}

	log.Println(h.meshTree.Tree)
	log.Println(ev)
	if exp, isExpander := ev.Src().(gwu.Expander); isExpander {
		if exp.Expanded() {
			if exp.Content() == nil {
				exp.SetToolTip("Click to expand/collapse")
				children, _ := getChildren(h.meshTree, h.db)
				if len(children) > 0 {
					log.Printf("MMMMMMMMMM %+v\n", ev)
					log.Printf("MMMMMMMMMM src %+v\n", ev.Src())
					log.Printf("MMMMMMMMMM type %+v\n", ev.Type())
					log.Printf("MMMMMMMMMM id %+v\n", ev.Src().Id())

					p := gwu.NewPanel()
					exp.SetContent(p)
					for i, _ := range children {
						child := children[i]
						log.Println("child:", child.Tree)
						numChildren := countChildren(child, h.db, true)
						linePanel := gwu.NewHorizontalPanel()
						p.AddVSpace(5)
						if numChildren > 0 {
							newe := gwu.NewExpander()
							newe.SetHeader(linePanel)

							titleLabel := gwu.NewHtml("<b>" + child.DescriptorName + "</b> [" + findLabel(child) + "]")
							linePanel.Add(titleLabel)
							linePanel.AddHSpace(24)
							linePanel.Add(gwu.NewLabel(strconv.FormatInt(numChildren, 10) + " " + makeChildWord(numChildren)))
							linePanel.AddHSpace(30)
							link := gwu.NewLink(child.DescriptorUI, "https://meshb-prev.nlm.nih.gov/#/record/ui?name="+child.DescriptorName)
							link.SetToolTip("NCBI MeSH Descriptor Record: " + child.DescriptorName)
							link.SetAttr("foo", "bar")
							linePanel.Add(link)

							newLinkHandler := &linkClickHandler{}
							link.AddEHandler(newLinkHandler, gwu.ETypeClick)

							meh := makeHandler(child, h.db)
							newe.AddEHandler(meh, gwu.ETypeStateChange)
							p.Add(newe)

						} else {

							p.Add(linePanel)
							link := gwu.NewLink(child.DescriptorUI, "https://meshb-prev.nlm.nih.gov/#/record/ui?name="+child.DescriptorName)
							link.SetToolTip("NCBI MeSH Descriptor Record: " + child.DescriptorName)
							p.Add(link)
							l := gwu.NewHtml("<b>" + child.DescriptorName + "</b> [" + findLabel(child) + "]")
							//l := gwu.NewLabel(findLabel(child) + ": " + child.DescriptorName)
							linePanel.Add(l)
							linePanel.AddHSpace(24)

							linePanel.Add(link)

						}
					}

				}
			}
			ev.MarkDirty(exp)
		}
	}
}

func makeChildWord(n int64) string {
	if n == 1 {
		return "descendent"
	} else {
		return "descendents"
	}
}

func main() {
	log.SetOutput(ioutil.Discard)
	handlerMap = make(map[int64]*meshExpanderHandler)
	db, err := dbOpen("mesh2016_sqlite3.db")
	if err != nil {
		log.Fatal(err)
	}
	//win := makeWin(db)

	// Create and start a GUI server (omitting error check)
	// server := gwu.NewServer("guitest", "localhost:8081")
	// server.SetText("Test GUI App")
	// server.AddWin(win)
	// server.Start("") // Also opens windows list in browser

	server := gwu.NewServer("guitest", "localhost:8081")
	server.AddSessCreatorName("main", "Login Window")
	server.AddSHandler(sessHandler{db: db})
	server.Start("") // Also opens windows list in browser
}

func makeWin(db *gorm.DB) gwu.Window {
	// Create and build a window
	win := gwu.NewWindow("main", "Test GUI Window")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HALeft)
	win.SetCellPadding(2)

	p := gwu.NewPanel()
	win.Add(p)
	l := gwu.NewLabel("PubMeSHicator")
	l.Style().SetColor(gwu.ClrGreen)
	p.Add(l)
	p.AddVSpace(5)

	//numChildren, topLevel, err := getTopLevel(db)
	_, topLevel, err := getTopLevel(db)
	if err != nil {
		log.Fatal(err)
	}
	for i, t := range topLevel {
		log.Println(i, t)
	}

	for _, t := range topLevel {
		numChildren := countChildren(t, db, true)

		meh := makeHandler(t, db)
		e := gwu.NewExpander()
		//e.SetHeader(gwu.NewLabel(*t.T0 + " " + t.DescriptorName + " [" + strconv.FormatInt(numChildren[0], 10) + "]"))
		//e.SetHeader(gwu.NewLabel(*t.T0 + " " + t.DescriptorName + " " + strconv.FormatInt(numChildren, 10)))
		e.SetHeader(gwu.NewLabel(findLabel(t) + " " + t.DescriptorName + " " + strconv.FormatInt(numChildren, 10) + " descendents"))

		e.AddEHandler(meh, gwu.ETypeStateChange)
		p.Add(e)
		p.AddVSpace(15)
	}

	return win
}

type sessHandler struct {
	db *gorm.DB
}

func (h sessHandler) Created(s gwu.Session) {
	win := makeWin(h.db)
	s.AddWin(win)
}
func (h sessHandler) Removed(s gwu.Session) {}

func getTopLevel(db *gorm.DB) ([]int64, []*MeshTree, error) {
	//var mt []MeshTree

	//db.Where("top is not null and t0 is null").Find(&mt)

	//return mt, nil
	return getLevel(0, db)
}

func getChildren(mt *MeshTree, db *gorm.DB) ([]*MeshTree, error) {
	q := getChildrenQuery(mt, false)

	log.Println("**************************", mt.Depth, q)
	var mtChildren []*MeshTree
	db.Where(q).Find(&mtChildren)
	var count int64
	db.Model(&MeshTree{}).Where(q).Count(&count)
	log.Println(count)
	q = getChildrenQuery(mt, true)
	log.Println("******* all", q)
	db.Model(&MeshTree{}).Where(q).Count(&count)
	log.Println(count)
	return mtChildren, nil
}

func countChildren(mt *MeshTree, db *gorm.DB, allChildren bool) int64 {
	var count int64

	q := getChildrenQuery(mt, allChildren)
	db.Model(&MeshTree{}).Where(q).Count(&count)
	return count
}

func getLevel(level int, db *gorm.DB) ([]int64, []*MeshTree, error) {
	var mt []*MeshTree
	q := "T0 is not null"
	for i := 1; i < level+2; i++ {
		q += " AND t" + strconv.Itoa(i) + " is null"
	}
	log.Println(q)
	db.Where(q).Find(&mt)

	var tmp []int64
	return tmp, mt, nil
}

func whichDescendents(all bool, l1, l2 string) string {
	if all {
		return " AND " + l1 + " IS NOT NULL"
	} else {
		r := " AND " + l1 + " IS NOT NULL"
		if l1 != "" {
			r += " AND " + l2 + " IS NULL"
		}
		return r
	}
}

const T0 = "T0"
const T1 = "T1"
const T2 = "T2"
const T3 = "T3"
const T4 = "T4"
const T5 = "T5"
const T6 = "T6"
const T7 = "T7"
const T8 = "T8"
const T9 = "T9"
const T10 = "T10"
const T11 = "T11"
const T12 = "T12"

func getChildrenQuery(mt *MeshTree, allDescendents bool) string {
	q := ""

	switch mt.Depth {
	case 0:
		q = T0 + "=\"" + *mt.T0 + "\"" + whichDescendents(allDescendents, T1, T2)
	case 1:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\"" + whichDescendents(allDescendents, T2, T3)
	case 2:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\" and " + T2 + "=\"" + *mt.T2 + "\"" + whichDescendents(allDescendents, T3, T4)
	case 3:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\" and " + T2 + "=\"" + *mt.T2 + "\" and " + T3 + "=\"" + *mt.T3 + "\"" + whichDescendents(allDescendents, T4, T5)
	case 4:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\" and " + T2 + "=\"" + *mt.T2 + "\" and " + T3 + "=\"" + *mt.T3 + "\" and " + T4 + "=\"" + *mt.T4 + "\"" + whichDescendents(allDescendents, T5, T6)

	case 5:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\" and " + T2 + "=\"" + *mt.T2 + "\" and " + T3 + "=\"" + *mt.T3 + "\" and " + T4 + "=\"" + *mt.T4 + "\" and " + T5 + "=\"" + *mt.T5 + "\"" + whichDescendents(allDescendents, T6, T7)
	case 6:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\" and " + T2 + "=\"" + *mt.T2 + "\" and " + T3 + "=\"" + *mt.T3 + "\" and " + T4 + "=\"" + *mt.T4 + "\" and " + T5 + "=\"" + *mt.T5 + "\" and " + T6 + "=\"" + *mt.T6 + "\"" + whichDescendents(allDescendents, T7, T8)
	case 7:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\" and " + T2 + "=\"" + *mt.T2 + "\" and " + T3 + "=\"" + *mt.T3 + "\" and " + T4 + "=\"" + *mt.T4 + "\" and " + T5 + "=\"" + *mt.T5 + "\" and " + T6 + "=\"" + *mt.T6 + "\" and " + T7 + "=\"" + *mt.T7 + "\"" + whichDescendents(allDescendents, T8, T9)
	case 8:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\" and " + T2 + "=\"" + *mt.T2 + "\" and " + T3 + "=\"" + *mt.T3 + "\" and " + T4 + "=\"" + *mt.T4 + "\" and " + T5 + "=\"" + *mt.T5 + "\" and " + T6 + "=\"" + *mt.T6 + "\" and " + T7 + "=\"" + *mt.T7 + "\" and " + T8 + "=\"" + *mt.T8 + "\"" + whichDescendents(allDescendents, T9, T10)

	case 9:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\" and " + T2 + "=\"" + *mt.T2 + "\" and " + T3 + "=\"" + *mt.T3 + "\" and " + T4 + "=\"" + *mt.T4 + "\" and " + T5 + "=\"" + *mt.T5 + "\" and " + T6 + "=\"" + *mt.T6 + "\" and " + T7 + "=\"" + *mt.T7 + "\" and " + T8 + "=\"" + *mt.T8 + "\" and T9=\"" + *mt.T9 + "\"" + whichDescendents(allDescendents, T10, T11)

	case 10:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\" and " + T2 + "=\"" + *mt.T2 + "\" and " + T3 + "=\"" + *mt.T3 + "\" and " + T4 + "=\"" + *mt.T4 + "\" and " + T5 + "=\"" + *mt.T5 + "\" and " + T6 + "=\"" + *mt.T6 + "\" and " + T7 + "=\"" + *mt.T7 + "\" and " + T8 + "=\"" + *mt.T8 + "\" and T9=\"" + *mt.T9 + "\" and T10=\"" + *mt.T10 + "\"" + whichDescendents(allDescendents, T11, T12)

	case 11:
		q = T0 + "=\"" + *mt.T0 + "\" and " + T1 + "=\"" + *mt.T1 + "\" and " + T2 + "=\"" + *mt.T2 + "\" and " + T3 + "=\"" + *mt.T3 + "\" and " + T4 + "=\"" + *mt.T4 + "\" and " + T5 + "=\"" + *mt.T5 + "\" and " + T6 + "=\"" + *mt.T6 + "\" and " + T7 + "=\"" + *mt.T7 + "\" and " + T8 + "=\"" + *mt.T8 + "\" and T9=\"" + *mt.T9 + "\" and T10=\"" + *mt.T10 + "\" and T11=\"" + *mt.T11 + "\"" + whichDescendents(allDescendents, T12, "")

	}

	return q

}
func findLabel(t *MeshTree) string {
	if t.T12 != nil {
		return *t.T0 + "." + *t.T1 + "." + *t.T2 + "." + *t.T3 + "." + *t.T4 + "." + *t.T5 + "." + *t.T6 + *t.T7 + "." + *t.T8 + "." + *t.T9 + "." + *t.T10 + "." + *t.T11 + "." + *t.T12
	} else {
		if t.T11 != nil {
			return *t.T0 + "." + *t.T1 + "." + *t.T2 + "." + *t.T3 + "." + *t.T4 + "." + *t.T5 + "." + *t.T6 + *t.T7 + "." + *t.T8 + "." + *t.T9 + "." + *t.T10 + "." + *t.T11
		} else {
			if t.T10 != nil {
				return *t.T0 + "." + *t.T1 + "." + *t.T2 + "." + *t.T3 + "." + *t.T4 + "." + *t.T5 + "." + *t.T6 + *t.T7 + "." + *t.T8 + "." + *t.T9 + "." + *t.T10
			} else {
				if t.T9 != nil {
					return *t.T0 + "." + *t.T1 + "." + *t.T2 + "." + *t.T3 + "." + *t.T4 + "." + *t.T5 + "." + *t.T6 + *t.T7 + "." + *t.T8 + "." + *t.T9
				} else {
					if t.T8 != nil {
						return *t.T0 + "." + *t.T1 + "." + *t.T2 + "." + *t.T3 + "." + *t.T4 + "." + *t.T5 + "." + *t.T6 + *t.T7 + "." + *t.T8
					} else {
						if t.T7 != nil {
							return *t.T0 + "." + *t.T1 + "." + *t.T2 + "." + *t.T3 + "." + *t.T4 + "." + *t.T5 + "." + *t.T6 + *t.T7
						} else {
							if t.T6 != nil {
								return *t.T0 + "." + *t.T1 + "." + *t.T2 + "." + *t.T3 + "." + *t.T4 + "." + *t.T5 + "." + *t.T6
							} else {
								if t.T5 != nil {
									return *t.T0 + "." + *t.T1 + "." + *t.T2 + "." + *t.T3 + "." + *t.T4 + "." + *t.T5
								} else {
									if t.T4 != nil {
										return *t.T0 + "." + *t.T1 + "." + *t.T2 + "." + *t.T3 + "." + *t.T4
									} else {
										if t.T3 != nil {
											return *t.T0 + "." + *t.T1 + "." + *t.T2 + "." + *t.T3
										} else {
											if t.T2 != nil {
												return *t.T0 + "." + *t.T1 + "." + *t.T2
											} else {
												if t.T1 != nil {
													return *t.T0 + "." + *t.T1
												} else {
													return *t.T0
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
