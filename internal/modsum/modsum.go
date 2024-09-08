package modsum

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/sumdb/dirhash"
)

func SumModDir(dir, modPath, modVersion string) (sum, goModSum string, err error) {
	prefix := modPath + "@" + modVersion

	files, err := dirFiles(dir, prefix)
	if err != nil {
		return
	}

	sum, err = dirhash.DefaultHash(files, func(s string) (io.ReadCloser, error) {
		return os.Open(filepath.Join(dir, strings.TrimPrefix(s, prefix)))
	})
	if err != nil {
		return
	}
	sum = modPath + " " + modVersion + " " + sum

	goModSum, err = dirhash.Hash1([]string{"go.mod"}, func(s string) (io.ReadCloser, error) {
		return os.Open(filepath.Join(dir, s))
	})
	goModSum = modPath + " " + modVersion + "/go.mod" + " " + goModSum

	return
}

func withGoMod(p string) bool {
	info, err := os.Stat(filepath.Join(p, "go.mod"))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// dirFiles returns the list of files in the tree rooted at dir,
// replacing the directory name dir with prefix in each name.
// The resulting names always use forward slashes.
func dirFiles(dir, prefix string) ([]string, error) {
	var files []string
	dir = filepath.Clean(dir)
	err := filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if filepath.Clean(file) == dir {
				return nil
			}

			if info.Name() == ".git" {
				return filepath.SkipDir
			}

			if withGoMod(file) {
				return filepath.SkipDir
			}

			return nil
		} else if file == dir {
			return fmt.Errorf("%s is not a directory", dir)
		}

		rel := file
		if dir != "." {
			rel = file[len(dir)+1:]
		}
		f := filepath.Join(prefix, rel)
		files = append(files, filepath.ToSlash(f))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
