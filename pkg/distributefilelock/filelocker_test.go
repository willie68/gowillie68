package distributefilelock

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLocker(t *testing.T) {
	ast := assert.New(t)
	file := "../../testdata/file.me"
	err := os.WriteFile(file, []byte("This is my file"), os.ModePerm)
	ast.NoError(err)

	lock := New(file)
	ast.NotNil(lock)
	ast.False(lock.IsLocked())

	ok := lock.TryLock()
	ast.True(ok)
	if ok {
		defer lock.Unlock()
	}

	slock := New(file)
	ast.NotNil(slock)
	ast.True(slock.IsLocked())

	ok = slock.TryLock()
	ast.False(ok)

	lock.Unlock()

	ast.False(slock.IsLocked())

}
