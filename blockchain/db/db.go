package db

import (
	"github.com/YoonBaek/CryptoProject/utils"
	"github.com/boltdb/bolt"
)

const (
	dbName      = "blockchain.db"
	dataBucket  = "data"
	blockBucket = "blocks"
	cp          = "checkpoint"
)

var db *bolt.DB

// This function is sigleton pattern of loading blockchain db
// if block chain db has not been loaded to memory,
// load db to memory, and construct buckets,
func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleErr(err)
		db = dbPointer

		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blockBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blockBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockchain(data []byte) {
	err := db.Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte("checkpoint"), data)
		return err
	})
	utils.HandleErr(err)
}

func LoadBlock(hash string) (blockBytes []byte) {
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blockBucket))
		blockBytes = bucket.Get([]byte(hash))
		return nil
	})
	return
}

func LoadCheckPoint() (data []byte) {
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(cp))
		return nil
	})
	return
}

func Close() {
	DB().Close()
}
