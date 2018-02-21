package util

import (
	"fmt"

	bolt "github.com/coreos/bbolt"
)

// OpenDb open the database
func OpenDb() (db *bolt.DB) {
	// open db
	var err error
	db, err = bolt.Open("dbFile/coinmarketcap.db", 0600, nil)
	if err != nil {
		fmt.Println("Cannot open db: ", err.Error())
		return nil
	}
	fmt.Println("DB opened")
	return
}

// CloseDb close the database
func CloseDb(db *bolt.DB) {
	db.Close()
	fmt.Println("DB closed")
}
