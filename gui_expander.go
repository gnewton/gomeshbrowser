package main

import (
	//"io/ioutil"
	"log"
	"strconv"

	"github.com/icza/gowut/gwu"
	"github.com/jinzhu/gorm"
)

func getTopLevel(db *gorm.DB) ([]int64, []*MeshTree, error) {
	//var mt []MeshTree

	//db.Where("top is not null and t0 is null").Find(&mt)

	//return mt, nil
	return getLevel(0, db)
}

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
		//log.Println("Cache hit")
		return val
	}
	//log.Println("Cache miss")
	newHandler := &meshExpanderHandler{meshTree: meshTree, db: db}
	handlerMap[meshTree.ID] = newHandler

	return newHandler
}

func (h *linkClickHandler) HandleEvent(ev gwu.Event) {

}

func (h *meshExpanderHandler) HandleEvent(ev gwu.Event) {
	log.Printf("ZZZZZZZZZZ %+v\n", ev)
	log.Printf("ZZZZZZZZZZ Event%+v\n", ev.Session())
	if hrr, ok := ev.(gwu.HasRequestResponse); ok {
		req := hrr.Request()
		log.Println("Client addr:", req.RemoteAddr)
		log.Println(ev)

	}

	//log.Println(h.meshTree.Tree)
	log.Println(ev)
	if exp, isExpander := ev.Src().(gwu.Expander); isExpander {
		if exp.Expanded() {
			if exp.Content() == nil {
				exp.SetToolTip("Click to expand/collapse")
				children, _ := getChildren(h.meshTree, h.db)
				if len(children) > 0 {
					//log.Printf("MMMMMMMMMM %+v\n", ev)
					// log.Printf("MMMMMMMMMM src %+v\n", ev.Src())
					// log.Printf("MMMMMMMMMM type %+v\n", ev.Type())
					// log.Printf("MMMMMMMMMM id %+v\n", ev.Src().Id())

					p := gwu.NewPanel()
					exp.SetContent(p)
					for i, _ := range children {
						child := children[i]
						//log.Println("child:", child.Tree)
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
