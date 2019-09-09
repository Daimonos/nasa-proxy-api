package main

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

var DB DBClient = DBClient{}

type DBClient struct {
	client *bbolt.DB
}

func (db *DBClient) Open(filepath string) error {
	d, err := bbolt.Open(filepath, 0600, nil)
	d.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("NEOS"))
		return err
	})
	if err != nil {
		return err
	}
	db.client = d
	return nil
}

func (db *DBClient) Set(bucket, key string, val interface{}) error {
	db.client.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		bytes, err := GetBufferFromStruct(val)
		if err != nil {
			return err
		}
		bucket.Put([]byte(key), bytes)
		return nil
	})
	return nil
}

func (db *DBClient) Get(bucket, key string) ([]byte, error) {
	var val []byte
	err := db.client.View(func(tx *bbolt.Tx) error {
		var err error
		bucket := tx.Bucket([]byte(bucket))
		val = bucket.Get([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

func GetBufferFromStruct(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
