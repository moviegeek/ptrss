package parser

import (
	"fmt"
	"os"

	"github.com/mmcdole/gofeed"
)

const (
	//HDCRssURL rss link to retrieve hdc movies, you need to add passkey at the end
	HDCRssURL = "https://hdchina.org/torrentrss.php?rows=50&cat17=1&cat9=1&isize=1"

	hdcPasskeyEnvVar = "HDC_PASSKEY"

	//PutaoRssURL rss link to retrieve putao movies, no need to add passkey
	PutaoRssURL = "https://pt.sjtu.edu.cn/torrentrss.php?rows=50&cat401=1&cat402=1&cat403=1&sta1=1&sta3=1&isize=1"
)

//RssParser is the main parser for both hdc and putao
type RssParser struct {
	hdcPasskey string
	parser     *gofeed.Parser
}

//GetHDCFeedItems parse HDC rss url and return items
func (p *RssParser) GetHDCFeedItems() ([]*gofeed.Item, error) {
	hdcRssURLWithPasskey := fmt.Sprintf("%s&passkey=%s", HDCRssURL, p.hdcPasskey)

	return p.getFeedItems(hdcRssURLWithPasskey)
}

//GetPutaoFeedItems parse Putao rss url and return items
func (p *RssParser) GetPutaoFeedItems() ([]*gofeed.Item, error) {
	return p.getFeedItems(PutaoRssURL)
}

func (p *RssParser) getFeedItems(rssURL string) ([]*gofeed.Item, error) {
	items := []*gofeed.Item{}

	feed, err := p.parser.ParseURL(rssURL)
	if err != nil {
		return items, err
	}

	return feed.Items, nil
}

//New creates a new RssParser
func New() *RssParser {
	hdcPasskey := os.Getenv(hdcPasskeyEnvVar)
	if hdcPasskey == "" {
		panic("Need to set passkey for HDChina in Enviroment variable: " + hdcPasskeyEnvVar)
	}

	return &RssParser{hdcPasskey, gofeed.NewParser()}
}
