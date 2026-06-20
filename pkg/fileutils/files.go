package fileutils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
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

type Callback func(fileinfo fs.DirEntry) bool

func GetFiles(rootpath, prefix string, callback Callback) error {
	infos, err := os.ReadDir(rootpath)
	if err != nil {
		return err
	}
	for _, i := range infos {
		if prefix == "" || strings.HasPrefix(strings.ToLower(i.Name()), strings.ToLower(prefix)) {
			ok := callback(i)
			if !ok {
				return nil
			}
		}
	}
	return nil
}

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

// SanitizePathSegment strips characters that are unsafe in directory names.
func SanitizePathSegment(s string) string {
	s = strings.TrimSpace(s)
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	s = replacer.Replace(s)
	s = strings.Map(func(r rune) rune {
		if r < 32 {
			return -1
		}
		return r
	}, s)
	s = strings.TrimRight(s, " .")
	if s == "" || s == "." || s == ".." {
		return "unknown"
	}
	if IsWindowsReservedName(s) {
		return "_" + s
	}
	return s
}

// SanitizeFileName returns a sanitized filename, removing or replacing characters that are unsafe in filenames.
func SanitizeFileName(name string) string {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "" || base == "." || base == ".." {
		return ""
	}

	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	stem = SanitizePathSegment(stem)
	if stem == "" || stem == "unknown" {
		stem = "file"
	}

	ext = strings.Map(func(r rune) rune {
		if r < 32 {
			return -1
		}
		if strings.ContainsRune("<>:\\|?*\"", r) {
			return -1
		}
		return r
	}, ext)
	ext = strings.TrimRight(ext, " .")

	if IsWindowsReservedName(stem) {
		stem = "_" + stem
	}

	return stem + ext
}

// IsWindowsReservedName checks if a name is a reserved name in Windows (e.g., CON, PRN, AUX, NUL, COM1, LPT1, etc.)
func IsWindowsReservedName(s string) bool {
	upper := strings.ToUpper(strings.TrimSpace(s))
	if upper == "" {
		return false
	}
	reserved := map[string]struct{}{
		"CON":  {},
		"PRN":  {},
		"AUX":  {},
		"NUL":  {},
		"COM1": {},
		"COM2": {},
		"COM3": {},
		"COM4": {},
		"COM5": {},
		"COM6": {},
		"COM7": {},
		"COM8": {},
		"COM9": {},
		"LPT1": {},
		"LPT2": {},
		"LPT3": {},
		"LPT4": {},
		"LPT5": {},
		"LPT6": {},
		"LPT7": {},
		"LPT8": {},
		"LPT9": {},
	}
	_, ok := reserved[upper]
	return ok
}
