package crawler

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/MatchlockIsDAST/sphttpclient/client"
	"github.com/PuerkitoBio/goquery"
)

//Crawler クローラーを定義する
type Crawler interface {
	Crawl(u *url.URL) (getlist []string, postlist map[string][]string)
}

//New Crawlerを新規で生成する
func New() Crawler {
	jar, err := cookiejar.New(nil)
	errorHandling(err)
	c := http.Client{
		Jar:     jar,
		Timeout: time.Duration(20) * time.Second,
	}
	return &crawler{
		client: client.New(c),
	}
}

type crawler struct {
	client client.Client
}

//クローリングする
func (c *crawler) Crawl(u *url.URL) (getlist []string, postlist map[string][]string) {
	var (
		memo  = map[string]bool{}
		crawl func(u *url.URL)
	)
	crawl = func(u *url.URL) { //URLをメモして、メモされていないものを飛ばす
		if !memo[u.String()] {
			resp := c.get(u)
			for tag := range attr {
				c.fetchLinks(u, resp, tag)
			}
			//crawl(u)
		}
	}
	crawl(u)
	return getlist, postlist
}

//Getしてくる
func (c *crawler) get(u *url.URL) (doc *goquery.Document) {
	req, err := http.NewRequest("GET", u.String(), nil)
	errorHandling(err)
	req.Header.Add("user-agent", "Matchlock Crawler v0.1")
	resp, err := c.client.Do(req)
	errorHandling(err)
	defer resp.Body.Close()
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	errorHandling(err)
	return doc
}

var attr = map[string]string{"a": "href", "img": "src", "script": "src", "link": "href"}

//Linkの吐き出し
func (c *crawler) fetchLinks(u *url.URL, doc *goquery.Document, tag string) {
	doc.Find(tag).Each(func(i int, a *goquery.Selection) {
		link, _ := a.Attr(attr[tag])
		u2 := linkParse(u, link)
		fmt.Println(u2.String())
	})
}

func linkParse(u *url.URL, link string) *url.URL {
	if strings.HasPrefix(link, "/") {
		u.Path = link
	} else if strings.HasPrefix(link, ".") {
		u.Path += link
	} else if strings.HasPrefix(link, "#") {
		u.Fragment = link
	} else {
		u, _ = url.Parse(link)
	}
	return u
}

func errorHandling(err error) {
	if err != nil {
		panic(err)
	}
}
