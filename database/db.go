package database

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open("data.db", 0600, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("data"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	db.Close()
}

func Get(key string) uint64 {
	var value uint64
	value = 0
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("data"))
		buf := bucket.Get([]byte(key))
		if buf != nil {
			value = binary.BigEndian.Uint64(buf[8:])
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return value
}

func Set(key string, value uint64) (uint64, error) {
	var old_value uint64
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("data"))
		timestamp := time.Now().YearDay()
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(timestamp))
		byte_value := make([]byte, 8)
		binary.BigEndian.PutUint64(byte_value, value)
		buf = append(buf, byte_value...)
		old_value_byte := bucket.Get([]byte(key))
		if old_value_byte != nil {
			old_value = binary.BigEndian.Uint64(old_value_byte[8:])
		}
		err := bucket.Put([]byte(key), buf)
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, err
	}
	return old_value, nil
}

func CountKeys() uint64 {
	var count uint64
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("data"))
		c := bucket.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			count++
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return count
}

func Monitor() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		err := db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("data"))
			c := bucket.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
                date := binary.BigEndian.Uint64(v[:8])
				now := uint64(time.Now().YearDay())
				fmt.Println(date, now)
                if now - date > 90 {
                    bucket.Delete(k)
                }
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
