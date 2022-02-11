package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

const DB = "todo.db"
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
		key := []byte(task.Name)
		fmt.Println(key)
		return bucket.Put(key, jsonb)
	})
}

func RemoveTask(nb int) {
	tasks := ReadTasks()
	if nb > len(tasks) {
		fmt.Println("No task #", nb)
	}

	task := tasks[nb-1]

	inBucket(func(bucket *bolt.Bucket) error {
		key := []byte(task.Name)
		return bucket.Delete(key)
	})

	task = tasks[nb-1]
}

func getDB() *bolt.DB {
	opts := &bolt.Options{Timeout: 1 * time.Second}
	db, err := bolt.Open(DB, 0600, opts)
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
