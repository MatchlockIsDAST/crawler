package sitemap

//SiteMap クローリングしたLinkを整理し利用するインターフェースを提供します
type SiteMap interface {
	AddChild(perURL string, childes SiteMap)
	AppendLinks(links map[string][]string)
	FetchLinks() map[string][]string
	FetchCrawlPage() string
	FetchChilde(keyURL string) SiteMap
}

type sitemap struct {
	CrawlPage string
	Links     map[string][]string
	Childs    map[string]SiteMap
}

//New 新規にSitemapを生成します
func New(crawlPage string, links map[string][]string) SiteMap {
	return &sitemap{
		CrawlPage: crawlPage,
		Links:     links,
		Childs:    map[string]SiteMap{},
	}
}

//Childを追加します
func (s *sitemap) AddChild(perURL string, childes SiteMap) {
	s.Childs[perURL] = childes
}

//Linksを追加します
func (s *sitemap) AppendLinks(links map[string][]string) {
	for key, value := range links {
		if len(s.Links[key]) > 0 {
			s.Links[key] = append(s.Links[key], value...)
		} else {
			s.Links[key] = value
		}
	}
}

//crawlしたLinkを提供します
func (s *sitemap) FetchLinks() map[string][]string { return s.Links }

//このSiteMapがどのページをcrawlしたかを確認できます
func (s *sitemap) FetchCrawlPage() string { return s.CrawlPage }

//Childを返します
func (s *sitemap) FetchChilde(keyURL string) SiteMap { return s.Childs[keyURL] }
