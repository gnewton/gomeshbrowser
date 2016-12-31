package main

import (
	//"io/ioutil"
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
	if exp, isExpander := ev.Src().(gwu.Expander); isExpander { // We clicked on an expander
		if exp.Expanded() { // We just expanded this expander
			if exp.Content() == nil { // No content in expander: we need to populate it from the DB
				populateNodeContent(exp, h.meshTree, h.db)
			}
			ev.MarkDirty(exp)
		}
	}
}

func populateNodeContent(exp gwu.Expander, meshTree *MeshTree, db *gorm.DB) {
	exp.SetToolTip("Click to expand/collapse")
	children, _ := getChildren(meshTree, db)
	if len(children) > 0 {
		p := gwu.NewPanel()
		exp.SetContent(p)
		for i, _ := range children {
			child := children[i]
			//log.Println("child:", child.Tree)
			numChildren := countChildren(child, db, true)

			linePanel := gwu.NewHorizontalPanel()
			p.AddVSpace(5)
			if numChildren > 0 {

				newExpander := gwu.NewExpander()
				newExpander.SetHeader(linePanel)
				makeExpanderContents(linePanel, child, numChildren)
				meh := makeHandler(child, db)
				newExpander.AddEHandler(meh, gwu.ETypeStateChange)
				p.Add(newExpander)

			} else {
				//p.Add(link)
				makeLeaf(linePanel, child)
				p.Add(linePanel)

			}
		}

	}
}

func makeLeaf(linePanel gwu.Panel, child *MeshTree) {

	l := gwu.NewHtml("<b>" + child.DescriptorName + "</b> [" + findLabel(child) + "]")
	//l := gwu.NewLabel(findLabel(child) + ": " + child.DescriptorName)
	linePanel.Add(l)
	linePanel.AddHSpace(24)
	link := gwu.NewLink(child.DescriptorUI, "https://meshb-prev.nlm.nih.gov/#/record/ui?name="+child.DescriptorName)
	link.SetToolTip("NCBI MeSH Descriptor Record: " + child.DescriptorName)
	linePanel.Add(link)
}

func makeExpanderContents(linePanel gwu.Panel, child *MeshTree, numChildren int64) {
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

}
