package cmd

import (
	"encoding/binary"
	"encoding/json"
	"hash/fnv"
	"log"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-homedir"
)

const BUCKET_NAME = "todo-bucket"

func ReadTasks() []Task {
	tasks := []Task{}

	inBucket(func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			var task Task
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}

			tasks = append(tasks, task)
			return nil
		})
	})

	return tasks
}

func AddTask(task Task) {
	jsonb, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	inBucket(func(bucket *bolt.Bucket) error {
		key := getKey(task.Name)
		return bucket.Put(key, jsonb)
	})
}

func RemoveTask(nb int) {
	task := getTask(nb)

	inBucket(func(bucket *bolt.Bucket) error {
		key := getKey(task.Name)
		return bucket.Delete(key)
	})
}

func CompleteTask(nb int) {
	task := getTask(nb)
	task.complete()
	AddTask(task)
}

func getTask(nb int) Task {
	tasks := ReadTasks()
	nb--
	if nb >= len(tasks) || nb < 0 {
		log.Fatalf("No task #%d", nb)
	}
	return tasks[nb]
}

func getDB() *bolt.DB {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	opts := &bolt.Options{Timeout: 1 * time.Second}
	dbPath := filepath.Join(home, ".todo.rc")
	db, err := bolt.Open(dbPath, 0600, opts)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func inBucket(f func(bucket *bolt.Bucket) error) {
	db := getDB()
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		if err != nil {
			return err
		}

		return f(bucket)
	})

	if err != nil {
		log.Fatal(err)
	}
}

func getKey(str string) []byte {
	sum := getHash(str)
	return toBytes(sum)
}

func getHash(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

func toBytes(i uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return b
}
