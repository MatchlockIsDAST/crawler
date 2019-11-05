package crawler

import (
	"net/url"
	"testing"
)

func TestCrawler(t *testing.T) {
	u, _ := url.Parse("https://qiita.com/A_zara")
	c := New()
	c.Crawl(u)
	t.Error()
}
