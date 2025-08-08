package distributefilelock

import (
	"fmt"
	"os"
	"time"

	utils "github.com/willie68/gowillie68/pkg"
)

type DistributedFileLock interface {
	Lock()                // waits until it can aquier a lock
	Unlock()              // unlock the given file
	TryLock() bool        // lock the file if possible, other wise return false
	IsLocked() bool       // Checks if a file is locked
	InstanceLocked() bool // Checks if a file is locked by this instance
}

type distributedFileLock struct {
	filename string
	lockname string
	locked   bool
}

func New(file string) DistributedFileLock {
	d := distributedFileLock{
		filename: file,
		lockname: file + ".lock",
		locked:   false,
	}

	return &d
}

func (d *distributedFileLock) Lock() {
	for {
		if d.TryLock() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (d *distributedFileLock) Unlock() {
	if !d.locked {
		return
	}
	if utils.FileExists(d.lockname) {
		_ = os.Remove(d.lockname)
		d.locked = false
	}
}

func (d *distributedFileLock) IsLocked() bool {
	if d.locked {
		return true
	}
	return utils.FileExists(d.lockname)
}

func (d *distributedFileLock) InstanceLocked() bool {
	return d.locked
}

func (d *distributedFileLock) TryLock() bool {
	if d.locked {
		return false
	}
	if !utils.FileExists(d.lockname) {
		f, err := os.Create(d.lockname)
		if err != nil {
			return false
		}
		defer f.Close()
		fmt.Fprintf(f, "%d", time.Now().Unix())
		d.locked = true
		return true
	}
	return false
}
