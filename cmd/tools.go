package cmd

import (
	"path/filepath"
	"os"
	"strings"
)

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}
