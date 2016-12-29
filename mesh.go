package main

import (
	//	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type MeshTree struct {
	ID             int64
	DescriptorUI   string `sql:"size:16"`
	DescriptorName string
	Tree           string
	Depth          int
	T0             *string `sql:"size:1"`
	T1             *string `sql:"size:3"`
	T2             *string `sql:"size:3"`
	T3             *string `sql:"size:3"`
	T4             *string `sql:"size:3"`
	T5             *string `sql:"size:3"`
	T6             *string `sql:"size:3"`
	T7             *string `sql:"size:3"`
	T8             *string `sql:"size:3"`
	T9             *string `sql:"size:3"`
	T10            *string `sql:"size:3"`
	T11            *string `sql:"size:3"`
	T12            *string `sql:"size:3"`
}

func dbOpen(dbFileName string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", dbFileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Opening db file: ", dbFileName)

	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	return db, nil
}
