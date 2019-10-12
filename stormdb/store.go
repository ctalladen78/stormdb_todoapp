package store

import (
	"log"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/json"
	"github.com/asdine/storm/q"
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

// user model
type User struct {
	ID     string ``
	Status []byte
	Name   []byte
}

// task model
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

func (s *Store) CreateUser(data *User) error {
	key := shortuuid.New()
	data.ID = key
	return s.Db.Save(data)
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

// Select a list of records that match a list of matchers. Doesn't use indexes.
// Select(matchers ...q.Matcher) Query
func (s *Store) FindOne(key string) (*Task, error) {
	output := &Task{}

	// dbOne("index", val, output)
	err := s.Db.One("ID", key, output)
	if err != nil {
		return nil, nil
	}
	return nil, nil
}

func (s *Store) DeleteTask(bucket string, key string) error {
	return nil
}

// filter sort transaction
func (s *Store) FilterTasksByStartDate() ([]*Task, error) {
	return nil, nil
}

// filter transaction
func (s *Store) GetDoneTasks() ([]*Task, error) {

	return nil, nil
}

// flutter client api call
// Future<List<UserData>> getUserListByCity({placeId: String}) async {
// get selected city placeId
// get users where currentLocation : City.placeId
// var queryParamString = "?where={\"currentLocation\":{\"\$inQuery\":{\"where\":{\"placeId\":\"$placeId\"},\"className\":\"City\"}}}";

// join relational transaction
// relational query between different objects
// task creator status: busy, available
// parse GET Comments 'where={"related_post":{"$inQuery":{"where":{"post_field":{"$equals":"abc123"}},"className":"Post"}}}'
func GetTasksWhereCreatorStatus(userdb string, taskdb string, status string) ([]*Task, error) {
	udb, err := storm.Open(userdb, storm.Codec(json.Codec))
	tdb, err := storm.Open(taskdb, storm.Codec(json.Codec))
	if err != nil {
		return nil, err
	}
	res := []*Task{}
	// TODO make nested queries
	// OPTION 1: client makes 2 queries, one to get filtered dataset,
	// then in client side map each item a concurrent function make parallel filter transaction
	m1 := q.Matcher(q.Eq("Status", status))
	// bolt.Select(...) doesnt use index
	q1 := userdb.Select(m1)
	q2 := taskdb.Select(m1)
	log.Print("QUERY ", q1)
	return res, nil
}

// map multiple items
// where={"title": "My post title", "likes": { "$gt": 100 }}'
func (s *Store) UpdateTasksAs(user string) error {
	return nil
}
