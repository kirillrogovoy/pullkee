package page

import (
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/kirillrogovoy/pullk/util"
)

type page struct {
	Response *http.Response
	Number   int
	Error    error
}

func (p page) getLastPageNumber() (int, bool) {
	links := p.Response.Header.Get("Link")
	if links == "" {
		return 0, false
	}

	return linksGetLastPageNumber(links)
}

// Pages is a sortable set of Pages
type Pages []page

func (a Pages) Len() int {
	return len(a)
}

func (a Pages) Swap(x, y int) {
	a[x], a[y] = a[y], a[x]
}

func (a Pages) Less(x, y int) bool {
	return a[x].Number < a[y].Number
}

// Fetcher is an interface for a fetcher of an entity with pages
type Fetcher interface {
	GetPage(page int) (*http.Response, error)
}

// FetchAll concurrently retrieves all pages of a particular entity using Fetcher
func FetchAll(fetcher Fetcher) (Pages, error) {
	var pages Pages

	firstPage := firstPage(fetcher)
	if err := firstPage.Error; err != nil {
		return nil, err
	}

	pages = append(pages, firstPage)
	lastPage := 1

	if l, hasAnyOtherPages := firstPage.getLastPageNumber(); hasAnyOtherPages {
		lastPage = l

		restOfPages, err := fetchConcurrently(fetcher, 2, lastPage)
		if err != nil {
			return nil, err
		}

		pages = append(pages, restOfPages...)
	}

	sort.Sort(pages)

	return pages, nil
}

func firstPage(f Fetcher) page {
	res, err := f.GetPage(1)
	return page{res, 1, err}
}

// fetches the last page number from the "Link" HTTP header which Github returns
func linksGetLastPageNumber(input string) (int, bool) {
	for _, link := range strings.Split(input, ",") {
		if strings.Contains(link, `rel="last"`) {
			page := regexp.MustCompile(`&page=(\d+)`).FindStringSubmatch(link)[1]
			if page == "" {
				log.Panicln("Couldn't find the last page in this link", link)
			}
			return util.ParseInt(page), true
		}
	}

	return 0, false
}

func fetchConcurrently(fetcher Fetcher, firstPage int, lastPage int) (Pages, error) {
	channel := make(chan page)

	for i := firstPage; i <= lastPage; i++ {
		go (func(i int) {
			res, err := fetcher.GetPage(i)
			channel <- page{res, i, err}
		})(i)
	}

	// wait for the results to come
	var pages Pages
	for i := firstPage; i <= lastPage; i++ {
		page := <-channel
		if err := page.Error; err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}
