package example

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/json"
	"github.com/asdine/storm/q"
)

type Db struct {
	DB *storm.DB
}

func NewDB(path string) (*Db, error) {
	d, err := storm.Open(path, storm.Codec(json.Codec))
	if err != nil {
		return nil, err
	}

	return &Db{DB: d}, nil
}

func (db *Db) Init(data interface{}) error {
	return db.DB.Init(data)
}

func (db *Db) All(bucket string, to interface{}) error {
	return db.DB.From(bucket).All(to)
}

func (db *Db) AllByIndex(bucket, key string, to interface{}) error {
	return db.DB.From(bucket).AllByIndex(key, to)
}

func (db *Db) Delete(bucket string, data interface{}) error {
	return db.DB.From(bucket).DeleteStruct(data)
}

func (db *Db) Select(bucket string, matchers ...q.Matcher) storm.Query {
	return db.DB.From(bucket).Select(matchers...)
}

func (db *Db) Update(bucket string, data interface{}) error {
	return db.DB.From(bucket).Update(data)
}
