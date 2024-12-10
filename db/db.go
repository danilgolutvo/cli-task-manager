package db

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

var db *bbolt.DB

const TaskBucket = "tasks"
const CompletedBucket = "completed"

func Init(dbPath string) {
	var err error
	db, err = bbolt.Open(dbPath, 0666, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Create the task bucket if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TaskBucket))
		if err != nil {
			return fmt.Errorf("failed to create tasks bucket: %w", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(CompletedBucket))
		if err != nil {
			return fmt.Errorf("failed to create completed tasks bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to create bucket:", err)
	}
}

func Add(task, buck string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		switch buck {
		case TaskBucket:
			bucket := tx.Bucket([]byte(TaskBucket))
			id, _ := bucket.NextSequence()
			return bucket.Put(itob(int(id)), []byte(task))
		default:
			bucket := tx.Bucket([]byte(CompletedBucket))
			id, _ := bucket.NextSequence()
			timestamp := time.Now().Unix()
			value := fmt.Sprintf("%d|%s", timestamp, task)
			return bucket.Put(itob(int(id)), []byte(value))
		}
	})
}
func List(buck string) ([]int, []string, error) {
	tasks := make(map[int]string)
	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(buck))
		if bucket == nil {
			return errors.New("task bucket does not exist")
		}

		// Fetch all tasks into the map
		return bucket.ForEach(func(k, v []byte) error {
			id := btoi(k) // Correctly decode the key
			tasks[id] = string(v)
			return nil
		})
	})

	if err != nil {
		return nil, nil, err
	}

	// Sort tasks by ID
	ids := make([]int, 0, len(tasks))
	for id := range tasks {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	// Collect tasks in sorted order
	sortedTasks := make([]string, 0, len(tasks))
	for _, id := range ids {
		sortedTasks = append(sortedTasks, tasks[id])
	}

	return ids, sortedTasks, nil
}
func Remove(id int, buck string, shouldRemoveAll bool) error {
	switch shouldRemoveAll {
	case true:
		return db.Update(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(CompletedBucket))
			return bucket.ForEach(func(k, v []byte) error {
				if err := bucket.Delete(k); err != nil {
					return err
				}
				return nil
			})
		})
	default:
		return db.Update(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(buck))
			return bucket.Delete(itob(id))
		})
	}
}

func Do(taskNum int) error {

	ids, sortedTasks, err := List(TaskBucket)

	dbid := ids[taskNum-1]
	err = Remove(dbid, TaskBucket, false)
	if err != nil {
		return err
	}
	if err := Add(sortedTasks[taskNum-1], CompletedBucket); err != nil {
		return err
	}
	return nil
}

func Completed() ([]int, []string, error) {
	tasks := make(map[int]string)
	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(CompletedBucket))
		if bucket == nil {
			return errors.New("task bucket does not exist")
		}

		// Fetch all tasks into the map
		var keysToDelete [][]byte

		err := bucket.ForEach(func(k, v []byte) error {
			parts := strings.SplitN(string(v), "|", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid value format: %s", string(v))
			}
			// Extract timestamp and task
			timestamp, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid timestamp: %s", parts[0])
			}

			task := parts[1]

			if isExpired(timestamp) {
				keysToDelete = append(keysToDelete, k)
			} else {
				id := btoi(k)
				tasks[id] = task
			}
			return nil
		})
		if err != nil {
			return err
		}
		// Delete expired tasks
		for _, key := range keysToDelete {
			if err := bucket.Delete(key); err != nil {
				return fmt.Errorf("failed to delete expired task: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	// Sort tasks by ID
	ids := make([]int, 0, len(tasks))
	for id := range tasks {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	// Collect tasks in sorted order
	sortedTasks := make([]string, 0, len(tasks))
	for _, id := range ids {
		sortedTasks = append(sortedTasks, tasks[id])
	}

	return ids, sortedTasks, nil
}

func Close() {
	if db != nil {
		db.Close()
	}
}
func isExpired(timestamp int64) bool {
	expirationTime := 24 * time.Hour
	return time.Now().After(time.Unix(timestamp, 0).Add(expirationTime))
}
func itob(v int) []byte {
	return []byte(string(rune(v)))
}

func btoi(b []byte) int {
	return int(rune(b[0]))
}
