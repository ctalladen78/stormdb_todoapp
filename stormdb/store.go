package store

import (
	"log"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/json"
	"github.com/lithammer/shortuuid"
	// "github.com/lithammer/shortuuid"
)

// "strings"
// "github.com/asdine/storm/q"
var taskBucket = []byte("users")

// https://godoc.org/github.com/asdine/storm#pkg-examples
// https://godoc.org/github.com/asdine/storm/q#Matcher
type Store struct {
	Db     *storm.DB
	Bucket []byte
}

type Task struct {
	ID     string ``
	Bucket string `storm:index`
	Status []byte
	Value  []byte
}

func NewDB(path string) (*Store, error) {
	d, err := storm.Open(path, storm.Codec(json.Codec))
	if err != nil {
		return nil, err
	}

	return &Store{Db: d}, nil
}

// initialize bucket before saving
// db starts with a main bucket
// possible to have sub-buckets or nested buckets
// func (s *Store) Init(bucket string, data *Task) error {
func (s *Store) Init(data interface{}) error {

	// TODO create bucket
	// return s.Db.Bolt.Update(func(tx *bbolt.Tx) error {
	// 	_, err := tx.CreateBucketIfNotExists([]byte(bucket))
	// 	return err
	// })
	return s.Db.Init(data)

}

func (s *Store) GetAll(dbPath string) ([]*Task, error) {
	log.Println("DB GET ", s.Db)
	res := []*Task{}
	err := s.Db.All(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Store) CreateTask(data *Task) error {
	log.Println("DB SAVE ", s.Db)
	// TODO setup tx to maintain process safety

	key := shortuuid.New()
	data.ID = key
	data.Status = []byte("open")

	// save to open bucket
	// return s.Db.From(bucket).Save(data)
	return s.Db.Save(data)

}

// move open task to completed bucket
func (s *Store) UpdateTask(key string, newValue string) error {
	currentTask := &Task{}
	currentTask.ID = key
	currentTask.Value = []byte(newValue)
	log.Printf("UPDATING %s", key)
	s.Db.Update(currentTask)
	return nil
}
func (s *Store) FindOne(idx int) (*Task, error) {
	return nil, nil

}

func (s *Store) DeleteTask(bucket string, key string) error {
	return nil
}
