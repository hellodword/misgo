package modsum

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func FindRepository(modPath, modVersion string) (repository string, err error) {
	fullURL := "https://pkg.go.dev/" + modPath + "@" + modVersion
	res, err := http.Get(fullURL)
	if err != nil {
		return
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	doc.Find(".UnitMeta-repo a").Each(func(i int, s *goquery.Selection) {
		repository, _ = s.Attr("href")
	})

	if repository == "" {
		err = fmt.Errorf("repository not found: %s", fullURL)
		return
	}

	return
}
