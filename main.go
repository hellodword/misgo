package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hellodword/misgo/internal/modsum"
	"golang.org/x/mod/sumdb"
)

func main() {
	dir := "tmp/mergo"
	gitURL := "https://github.com/imdario/mergo"
	modName, modVersion := "dario.cat/mergo", "v1.0.0"

	if _, err := os.Stat(dir); err != nil && errors.Is(err, os.ErrNotExist) {
		_, err := git.PlainClone(dir, false, &git.CloneOptions{
			URL:           gitURL,
			ReferenceName: plumbing.ReferenceName("refs/tags/" + modVersion),
			SingleBranch:  true,
			Depth:         1,
			Progress:      os.Stdout,
		})
		if err != nil {
			panic(err)
		}
	}

	// lookup from sum.golang.org
	ops := modsum.NewGoChecksumDatabaseClient()
	c := sumdb.NewClient(ops)
	records, err := c.Lookup(modName, modVersion)
	if err != nil {
		panic(err)
	}
	if len(records) != 1 {
		panic(records)
	}

	fmt.Println(records[0])

	sum, _, err := modsum.ModSum(dir, modName, modVersion)
	if err != nil {
		panic(err)
	}

	fmt.Println(sum)
	fmt.Println(sum == records[0])
}
