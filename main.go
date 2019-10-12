package main

import (
	"golang-projects/stormdb_todoapp/stormdb"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/asdine/storm/q"
)

var (
	db        store.Store
	storePath string
)

type Node struct {
	ID string
}

const (
	ALL      = "all"
	NEW      = "new"
	DELETE   = "delete"
	UPDATE   = "update"
	FIND     = "find"
	RUNTESTS = "test"
)

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	storePath = filepath.Join(dir, "/storm.db")
	log.Println(storePath)

	db, err := store.NewDB(storePath)
	if err != nil {
		log.Fatal("DB NEW error ", err)
	}

	// init with new bucket, initial user data
	err = db.Init(&store.Task{})
	err = db.Init(&store.User{})

	u1 := &store.User{Status: []byte("busy"), Name: []byte("sara")}
	u2 := &store.User{Status: []byte("busy"), Name: []byte("chris")}
	err = db.CreateUser(u1)
	err = db.CreateUser(u2)
	// err = db.CreateUser(&store.User{Status: []byte("busy"), Name: []byte("sara")})
	// err = db.CreateUser(&store.User{Status: []byte("open"), Name: []byte("chris")})

	if err != nil {
		log.Fatal("DB init error ", err)
	}

	input := os.Args[1]
	if input == "" {
		log.Fatal("pls enter a command")
	}
	log.Print("user input: ", string(input))

	switch input {
	case ALL:
		// Get All
		getAllResult, err := db.GetAll(storePath)
		if err != nil {
			log.Fatal("GET all error", err)
		}
		// first item should be updated
		log.Println("GET ALL")
		for i, s := range getAllResult {
			log.Printf("%d %s", i, s)
		}
		bucket := db.Db.Bucket() // should print bucket name
		log.Printf("bucket %s", bucket)
	case NEW:
		todoval := os.Args[2]
		if input == "" || todoval == "" {
			log.Fatal("pls enter a command")
		}
		t := &store.Task{Value: []byte(todoval)}
		// Create
		err = db.CreateTask(t)
		if err != nil {
			log.Fatal("Create error ", err)
		}
	case UPDATE:
		newval := os.Args[3]
		idx, err := strconv.Atoi(os.Args[2])
		if input == "" || newval == "" || err != nil {
			log.Fatal("pls enter a command")
		}
		getAllResult, err := db.GetAll(storePath)
		if len(getAllResult) > 1 {
			// r := rand.New(rand.NewSource(99)) // Update random item
			// rnum := r.Intn(len(getAllResult))
			// err = db.UpdateTask(getAllResult[rnum].ID, newval)
			err = db.UpdateTask(getAllResult[idx].ID, newval)
		}
		if err != nil {
			log.Fatal("Create error ", err)
		}
	case FIND: // find by list index
		idx, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("pls enter a command")
		}
		getAllResult, err := db.GetAll(storePath) // post transaction routine TODO in db routine
		// to := &store.Task{}
		to, err := db.FindOne(getAllResult[idx].ID)

		if err != nil || to != nil {
			log.Fatal("GET all error", err)
		}
		// foundTask, err := FindOne(idx)
		log.Printf("FIND ONE: %s", to.Value)
	case RUNTESTS:
		db.GetTasksWhereCreatorStatus("open") // relational
		db.GetDoneTasks()                     // filter
		db.UpdateTasksAs("done")              // map multiple items

	}

	// for serializing into json, convert []byte to string

	// Filter query using matchers
	var filters []q.Matcher
	var completed []*store.Task
	filters = append(filters, q.Eq("Status", "completed"))

	// query.Find(&result) executes query and puts results into result
	// err = db.Db.Select(filters...).Bucket("Entry").Find(index)
	q1 := db.Db.Select(filters...).Bucket("completed")
	if err = q1.Find(&completed); err != nil {
		return
	}

	// Select query using q ranges
	q2 := db.Db.Select(q.Lt("Date", time.Now()))
	if q2 != nil {
		return
	}
	q2.Find(&completed)

}
