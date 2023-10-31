package fetcher

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

	for i := 1; i < 50; i++ {
		url := fmt.Sprintf("%s%d", sourceUrl, i)

		wg.Add(1)

		go func(url string) {
			f.fetch(url, &wg)
		}(url)

	}

	wg.Wait()

	return nil
}

func (f *Fetcher) fetch(link string, wg *sync.WaitGroup) {
	defer wg.Done()

	response, err := http.Get(link)
	if err != nil {
		log.Printf("can't get responce from %s: %v", link, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("%s: status code error: %s\n", link, response.Status)
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
		thunbUrl, _ := s.Find(".tm-article-body").Find("img").Attr("src")
		t, _ := s.Find("time").Attr("title")

		u, err := url.JoinPath("https://habr.com/", link)
		if err != nil {
			log.Printf("[ERROR] cant't join url: %v", err)
		}

		article := db.Article{
			Title:        title,
			URL:          u,
			ThumbURL:     thunbUrl,
			Published_at: t,
		}

		f.Articles = append(f.Articles, article)

	})
}
