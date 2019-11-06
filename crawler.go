package crawler

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/MatchlockIsDAST/crawler/sitemap"
	"github.com/MatchlockIsDAST/sphttpclient/client"
	"github.com/PuerkitoBio/goquery"
)

//Crawler クローラーを定義する
type Crawler interface {
	Crawl(u *url.URL, depth int) (getSiteMap, postSiteMap sitemap.SiteMap)
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
func (c *crawler) Crawl(u *url.URL, depth int) (getSiteMap, postSiteMap sitemap.SiteMap) {
	var (
		memo  = map[string]bool{} //crawl前に確かめるmemo
		crawl func(u *url.URL, depth int) (getSiteMap, postSiteMap sitemap.SiteMap)
	)
	crawl = func(u *url.URL, depth int) (getSiteMap, postSiteMap sitemap.SiteMap) { //URLをメモして、メモされていないものを飛ばす
		if depth == 0 {
			return nil, nil
		}
		if !memo[u.String()] {
			memo[u.String()] = true
			resp := c.get(u)
			links := map[string][]string{}
			for tag := range attr {
				links[tag] = c.fetchLinks(u, resp, tag)
			}
			getSiteMap = sitemap.New(u.String(), links)
			for _, u := range links["a"] {
				up, _ := url.Parse(u)
				getSiteMapChild, _ := crawl(up, depth-1)
				if getSiteMapChild != nil {
					getSiteMap.AddChild(up.String(), getSiteMapChild)
				}
			}
		}
		return getSiteMap, postSiteMap
	}
	getSiteMap, postSiteMap = crawl(u, depth)
	return getSiteMap, postSiteMap
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
func (c *crawler) fetchLinks(u *url.URL, doc *goquery.Document, tag string) (links []string) {
	doc.Find(tag).Each(func(i int, a *goquery.Selection) {
		link, _ := a.Attr(attr[tag])
		u2 := linkParse(u, link)
		if u2.String() != "" {
			links = append(links, u2.String())
		}
	})
	return links
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
