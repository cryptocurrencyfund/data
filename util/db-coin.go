package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	bolt "github.com/coreos/bbolt"
	coinMarketCap "github.com/cryptocurrencyfund/go-coinmarketcap"
)

// UpdateCoin updates the coin in db
func UpdateCoin(db *bolt.DB, ts int64, c coinMarketCap.Coin) {
	coinType := strings.ToLower(c.Name)
	fmt.Printf("Updating coin: %s @%d\n", coinType, ts)
	coinBytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println("json marshall failed:", err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(coinType))
		if err != nil {
			fmt.Println("Creating bucket failed: " + coinType)
			return err
		}
		err = b.Put(Int64ToByteArr(ts), coinBytes)
		return err
	})
}

// GetCoin gets the coin in db
func GetCoin(db *bolt.DB, coinName string) map[int64]coinMarketCap.Coin {
	entries := make(map[int64]coinMarketCap.Coin)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(coinName))
		if b == nil {
			fmt.Println("Cannot find bucket: " + coinName)
			return errors.New("Cannot find bucket: " + coinName)
		}
		c := b.Cursor()
		counter := 0

		for k, v := c.First(); k != nil; k, v = c.Next() {
			// fmt.Printf("key=%s, value=%s\n", k, v)
			var data coinMarketCap.Coin
			err := json.Unmarshal(v, &data)
			if err != nil {
				fmt.Printf("Cannot unmarshall coin: %s\n", v)
				return errors.New("Cannot unmarshall")
			}

			entries[ByteArrToInt64(k)] = data
			counter++
		}

		fmt.Println("Total entries: ", counter)
		return nil
	})

	return entries
}

// ClearBucket delete everything and recreate bucket
func ClearBucket(db *bolt.DB, bucket string) {
	db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucket))
		if err != nil {
			fmt.Println("Error deleting bucket: "+bucket, err.Error())
		}
		tx.CreateBucketIfNotExists([]byte(bucket))

		return err
	})
}
