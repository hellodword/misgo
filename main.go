package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aquasecurity/trivy/pkg/dependency/parser/golang/sum"
	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hellodword/misgo/internal/modsum"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/sumdb"
)

func main() {
	log.SetLevel(log.InfoLevel)

	base := "/tmp/misgo"
	os.MkdirAll(base, 0755)

	modRoot := "."

	gomod := filepath.Join(modRoot, "go.mod")
	data, err := os.ReadFile(gomod)
	if err != nil {
		panic(err)
	}
	modInfo, err := modfile.Parse(gomod, data, nil)
	if err != nil {
		panic(err)
	}

	gosum := filepath.Join(modRoot, "go.sum")
	data, err = os.ReadFile(gosum)
	if err != nil {
		panic(err)
	}
	sumInfo, _, err := sum.NewParser().Parse(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	var deps = make(map[string]module.Version)

	for _, require := range modInfo.Require {
		deps[require.Mod.String()] = require.Mod
	}

	for _, replace := range modInfo.Replace {
		deps[replace.New.String()] = replace.New
	}

	for _, pkg := range sumInfo {
		if !strings.HasPrefix(pkg.Version, "v") {
			pkg.Version = "v" + pkg.Version
		}
		m := module.Version{Path: pkg.Name, Version: pkg.Version}
		deps[m.String()] = m
	}

	for _, dep := range deps {
		if strings.HasPrefix(dep.Path, "github.com/") ||
			strings.HasPrefix(dep.Path, "gitlab.com/") ||
			strings.HasPrefix(dep.Path, "golang.org/x/") {
			log.Debug("Skip", "mod", dep.String())
			continue
		}

		err := verify(base, dep.Path, dep.Version)
		if err != nil {
			log.Error("Failed", "mod", dep.String(), "err", err)
		} else {
			log.Info("Verified", "mod", dep.String())
		}
	}

}

func verify(base, modPath, modVersion string) error {

	dir := filepath.Join(base, filepath.Base(modPath))

	gitURL, err := modsum.FindRepository(modPath, modVersion)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dir); err != nil && errors.Is(err, os.ErrNotExist) {
		_, err := git.PlainClone(dir, false, &git.CloneOptions{
			URL:           gitURL,
			ReferenceName: plumbing.ReferenceName("refs/tags/" + modVersion),
			SingleBranch:  true,
			Depth:         1,
			Progress:      os.Stdout,
		})
		if err != nil {
			return err
		}
	}

	// lookup from sum.golang.org
	ops := modsum.NewGoChecksumDatabaseClient()
	c := sumdb.NewClient(ops)
	records, err := c.Lookup(modPath, modVersion)
	if err != nil {
		return err
	}
	if len(records) != 1 {
		return fmt.Errorf("too much records: %+v", records)
	}

	checksum, _, err := modsum.SumModDir(dir, modPath, modVersion)
	if err != nil {
		return err
	}

	if checksum != records[0] {
		return errors.New("checksum mismatch")
	}

	return nil
}
