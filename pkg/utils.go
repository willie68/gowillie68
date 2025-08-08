package utils

import (
	"fmt"
	"io"
	"slices"
	"strconv"
)

func ReadString(r io.Reader) (string, error) {
	bodyBytes, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("can't read body: %v", err)
	}
	return string(bodyBytes), nil
}

func ReadBytes(r io.Reader) ([]byte, error) {
	bodyBytes, err := io.ReadAll(r)
	if err != nil {
		return []byte{}, fmt.Errorf("can't read body: %v", err)
	}
	return bodyBytes, nil
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return false
	}
	return true
}

func Compact(s []string) []string {
	slices.Sort(s)
	return slices.Compact(s)
}
