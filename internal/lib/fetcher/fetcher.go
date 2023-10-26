package fetcher

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/yellowpuki/tg-bath-bot/internal/storage/db"
)

const sourceUrl = "https://habr.com/ru/hubs/go/articles/page"

type Fetcher struct {
	Articles      []db.Article
	fetchInterval time.Duration
}

func New(fetchInterval time.Duration) *Fetcher {
	articles := make([]db.Article, 0)

	return &Fetcher{
		Articles:      articles,
		fetchInterval: fetchInterval,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		}
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	var wg sync.WaitGroup

	for i := 1; i < 40; i++ {
		url := fmt.Sprintf("%s%d", sourceUrl, i)

		wg.Add(1)

		go func(url string) {
			f.fetch(url, &wg)
		}(url)

	}

	wg.Wait()

	return nil
}

func (f *Fetcher) fetch(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	response, err := http.Get(url)
	if err != nil {
		log.Printf("can't get responce from %s: %v", url, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("%s: status code error: %s\n", url, response.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Printf("can't parse response: %v", err)
		return
	}

	doc.Find(".tm-articles-list__item").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Find("h2").Find("span").Html()
		link, _ := s.Find("h2").Find("a").Attr("href")

		article := db.Article{
			Title: title,
			URL:   path.Join("https://habr.com/", link),
		}

		f.Articles = append(f.Articles, article)

	})
}
