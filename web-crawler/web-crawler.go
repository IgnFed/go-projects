package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

var cachedValues = cache{store: make(map[string]bool)}

func (c *cache) exist(url string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.store[url] {
		return true
	}
	c.store[url] = true
	return false
}

func recursiveCralw(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()
	if depth <= 0 || cachedValues.exist(url) {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wg.Add(1)
		go recursiveCralw(u, depth-1, fetcher, wg)
	}
	return
}

func Crawl(url string, depth int, fetcher Fetcher) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go recursiveCralw(url, depth, fetcher, wg)
	wg.Wait()

}

func WebCrawler() {
	Crawl("https://golang.org/", 4, fetcher)
}

type fakeFetcher map[string]fakeResult

type fakeResult struct {
	body string
	urls []string
}

type cache struct {
	mux   sync.Mutex
	store map[string]bool
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

var fetcher = fakeFetcher{
	"https://golang.org/": fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
