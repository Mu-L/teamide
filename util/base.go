package util

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	BaseDir string
)

func init() {
	var err error
	BaseDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	BaseDir, err = filepath.Abs(BaseDir)
	if err != nil {
		panic(err)
	}
	BaseDir = filepath.ToSlash(BaseDir)
	if !strings.HasSuffix(BaseDir, "/") {
		BaseDir += "/"
	}
}