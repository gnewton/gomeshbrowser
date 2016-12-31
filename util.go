package main

import (
	//	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

func getChildren(mt *MeshTree, db *gorm.DB) ([]*MeshTree, error) {
	var mtChildren []*MeshTree
	if mt.Depth == 12 {
		return mtChildren, nil
	}
	q := getChildrenQuery(mt, false)

	log.Println("**************************", mt.Depth, q)

	db.Where(q).Find(&mtChildren)
	// var count int64
	// db.Model(&MeshTree{}).Where(q).Count(&count)
	// log.Println(count)
	// q = getChildrenQuery(mt, true)
	// log.Println("******* all", q)
	// db.Model(&MeshTree{}).Where(q).Count(&count)
	// log.Println(count)
	return mtChildren, nil
}

func countChildren(mt *MeshTree, db *gorm.DB, allChildren bool) int64 {
	if mt.Depth == 12 {
		return 0
	}
	var count int64 = 0

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
	//log.Println(q)
	db.Where(q).Find(&mt)

	var tmp []int64
	return tmp, mt, nil
}

func whichDescendents(all bool, l1, l2 string) string {
	if all {
		if l1 != "" {
			return " AND " + l1 + " IS NOT NULL"
		} else {
			return ""
		}
	} else {
		r := " AND " + l1 + " IS NOT NULL"
		if l1 != "" && l1 != T12 {
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

	//log.Println("Depth=" + strconv.Itoa(mt.Depth))

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

func makeChildWord(n int64) string {
	if n == 1 {
		return DESCENDANT
	} else {
		return DESCENDANTS
	}
}
