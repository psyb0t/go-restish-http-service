package restishhttpservice

import (
	"fmt"
	"io/ioutil"
	"os"
)

type DB struct {
	path string
}

func NewDB(path string) *DB {
	db := &DB{
		path: path,
	}

	os.MkdirAll(db.path, 0755)

	return db
}

func (db *DB) SetPath(path string) {
	db.path = path
	os.MkdirAll(db.path, 0755)
}

func (db *DB) GetPath() string {
	return db.path
}

func (db *DB) MakeKeyPath(key string) string {
	return fmt.Sprintf("%s/%s", db.path, key)
}

func (db *DB) Set(key string, val []byte) error {
	return ioutil.WriteFile(db.MakeKeyPath(key), val, 0644)
}

func (db *DB) Get(key string) ([]byte, error) {
	return ioutil.ReadFile(db.MakeKeyPath(key))
}

func (db *DB) Delete(key string) error {
	return os.Remove(db.MakeKeyPath(key))
}
