package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/willie68/gowillie68/pkg/logging"
)

var (
	replaceChars = []string{"\"", "*", "/", ":", "<", ">", "?", "\\", "|", "."}
	removeChars  = []string{"\t", "\r", "\n"}
)

// FileExists checks if a file exsists
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else {
		return false
	}
}

// IsDir checks if a file exsists
func IsDir(filename string) bool {
	if f, err := os.Stat(filename); (err == nil) && (f.IsDir()) {
		return true
	} else {
		return false
	}
}

// FileNameWithoutExtension returning the filename without the extension
func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// ValidPathName return a valid path name, all non file chars will be changed to _
func ValidPathName(s string) string {
	s, _ = url.PathUnescape(s)
	for _, remove := range removeChars {
		s = strings.ReplaceAll(s, remove, "")
	}
	for _, replace := range replaceChars {
		s = strings.ReplaceAll(s, replace, "_")
	}
	return s
}

// HashFile build a sha256 hash of the file content
func HashFile(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		logging.Root.Errorf("%v", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		logging.Root.Errorf("%v", err)
	}
	return hex.EncodeToString(h.Sum(nil))
}
