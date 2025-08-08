package queuelist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	utils "github.com/willie68/gowillie68/pkg"
	"github.com/willie68/gowillie68/pkg/distributefilelock"
	"github.com/willie68/gowillie68/pkg/logging"
)

// QueueList simple queue backuped by a file, for smaller queues
type QueueList struct {
	name     string
	filename string
	list     []string
	slock    sync.Mutex
	log      logging.Logger
	dlock    distributefilelock.DistributedFileLock
}

// New  creates a new Queuelist in the folder path with a name (filename ist <name>.json)
func New(path string, name string) *QueueList {
	filename := filepath.Join(path, fmt.Sprintf("%s.json", name))
	q := QueueList{
		name:     name,
		filename: filename,
		list:     make([]string, 0),
		slock:    sync.Mutex{},
		dlock:    distributefilelock.New(filename),
		log:      *logging.New().WithName(fmt.Sprintf(`queue "%s"`, name)),
	}
	q.loadQueue()
	return &q
}

// Push puts a new value to the list, duplicates are not permitted
func (q *QueueList) Push(v string) bool {
	q.slock.Lock()
	defer q.slock.Unlock()
	if !slices.Contains(q.list, v) {
		q.list = append(q.list, v)
		q.saveQueue()
		return true
	}
	return false
}

// DirectPush puts a new value to the list, duplicates are not permitted, without saving
func (q *QueueList) DirectPush(v string) bool {
	q.slock.Lock()
	defer q.slock.Unlock()
	if !slices.Contains(q.list, v) {
		q.list = append(q.list, v)
		return true
	}
	return false
}

// Save saving the queue
func (q *QueueList) Save() {
	q.slock.Lock()
	defer q.slock.Unlock()
	q.saveQueue()
}

// Pop return the first entry from the list and deletes it
func (q *QueueList) Pop() (string, bool) {
	q.slock.Lock()
	defer q.slock.Unlock()
	if q.Size() == 0 {
		return "", false
	}
	v := q.list[0]
	q.list = slices.Delete(q.list, 0, 1)
	q.saveQueue()
	return v, true
}

// Peek return the first entry from the list
func (q *QueueList) Peek() (string, bool) {
	q.slock.Lock()
	defer q.slock.Unlock()
	if q.Size() == 0 {
		return "", false
	}
	v := q.list[0]
	return v, true
}

// Remove removes the entry from the queue
func (q *QueueList) Remove(v string) bool {
	q.slock.Lock()
	defer q.slock.Unlock()
	idx := slices.Index(q.list, v)
	if idx < 0 {
		return false
	}
	q.list = slices.Delete(q.list, idx, idx+1)
	q.saveQueue()
	return true
}

func (q *QueueList) saveQueue() {
	list, err := q.loadList()
	if err != nil {
		return
	}

	for _, x := range list {
		if !slices.Contains(q.list, x) {
			q.list = append(q.list, x)
		}
	}

	js, err := json.Marshal(q.list)
	if err != nil {
		q.log.Infof("can't serialise json: %v", err)
	}
	err = os.WriteFile(q.filename, js, os.ModePerm)
	if err != nil {
		q.log.Infof("can't write queue file: %v", err)
	}
}

func (q *QueueList) loadQueue() {
	list, err := q.loadList()
	if err != nil {
		return
	}
	q.list = list
}

func (q *QueueList) loadList() ([]string, error) {
	q.dlock.Lock()
	defer q.dlock.Unlock()

	list := make([]string, 0)
	if !utils.FileExists(q.filename) {
		return list, nil
	}

	fs, err := os.ReadFile(q.filename)
	if err != nil {
		q.log.Infof("can't write queue file: %v", err)
		return list, err
	}
	err = json.Unmarshal(fs, &list)
	if err != nil {
		q.log.Infof("can't deserialise json: %v", err)
		return list, err
	}
	return list, nil
}

// List return a copy slice of entries
func (q *QueueList) List() []string {
	tmp := make([]string, len(q.list))
	copy(tmp, q.list)
	return tmp
}

// Size return the length of this queue
func (q *QueueList) Size() int {
	return len(q.list)
}

func (q *QueueList) Clear() {
	q.slock.Lock()
	defer q.slock.Unlock()
	q.list = make([]string, 0)
}
