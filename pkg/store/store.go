package store

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/moviegeek/pt"
)

const (
	hdcSiteName     = "HDC"
	putaoSiteName   = "Putao"
	unknonwSiteName = "Unknown"
)

const contentTemplate = `
{{if .IMDBID}}
<p>IMDB: 
	<a href=https://www.imdb.com/title/{{.IMDBID}}>
		<span>{{.IMDBRating}} / {{.IMDBVotes}}</span>
	</a>
</p>
{{end}}
{{if .Poster}}
<div>
	<img alt="{{.Title}} Poster" title="{{.Title}}" src="{{.Poster}}"/>
</div>
{{end}}
<b/>
<p>Download:</p>
{{range .PTMedias}}
<div>
	<a href={{.Link}}>
		<span>{{.Site}} {{.MovieInfo.Source}} {{.MovieInfo.Resolution}} {{.MovieInfo.Group}} {{.MovieInfo.Size}}</span>
	</a>
</div>
{{end}}
`

//Store stores movie info
type Store struct {
	movies       []Movie
	indexByTitle map[string]*Movie
}

//AddFromFeedItem add a new movie info from an rss feed item
func (s *Store) AddFromFeedItem(item *gofeed.Item) {
	info := pt.ParseTitle(item.Title)

	if m, ok := s.indexByTitle[info.Title]; ok {
		s.addExsitMovieFromFeedItem(m, info, item)
	} else {
		s.addNewMovie(info, item)
	}
}

func (s *Store) addExsitMovieFromFeedItem(movie *Movie, info pt.MovieInfo, item *gofeed.Item) {
	media := PTMedia{
		MovieInfo: info,
		Site:      getSiteNameFromURL(item.Link),
		Link:      item.Link,
		SiteID:    getSiteIDFromURL(item.Link),
	}

	movie.PTMedias = append(movie.PTMedias, media)
}

func (s *Store) addNewMovie(info pt.MovieInfo, item *gofeed.Item) {
	movie := Movie{
		Title:     info.Title,
		Year:      info.Year,
		Updated:   item.Updated,
		Published: item.Published,
	}

	movie.PTMedias = []PTMedia{
		{
			MovieInfo: info,
			Site:      getSiteNameFromURL(item.Link),
			Link:      item.Link,
			SiteID:    getSiteIDFromURL(item.Link),
		},
	}

	s.movies = append(s.movies, movie)
	s.indexByTitle[movie.Title] = &s.movies[len(s.movies)-1]
}

//ToRss generate rss xml from the movies in this Store
func (s *Store) ToRss(w io.Writer) error {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "private torrent feeds",
		Description: "pt",
		Created:     now,
		Link:        &feeds.Link{Href: ""},
	}

	tmpl, err := template.New("content").Parse(contentTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse content template: %v", err)
	}

	for _, movie := range s.movies {
		item := movieToItem(&movie, tmpl)
		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, rss)
	return err
}

//Movies returns all existing movies in this store
func (s *Store) Movies() []Movie {
	return s.movies
}

func getSiteNameFromURL(url string) string {
	if strings.Contains(url, "pt.sjtu.edu.cn") {
		return putaoSiteName
	} else if strings.Contains(url, "hdchina.org") {
		return hdcSiteName
	}

	return unknonwSiteName
}

func getSiteIDFromURL(url string) string {
	id := ""
	i := strings.LastIndex(url, "id=")
	if i > -1 {
		id = url[i+3:]
		id = strings.TrimSpace(id)
	}
	return id
}

func movieToItem(movie *Movie, tmpl *template.Template) *feeds.Item {
	item := &feeds.Item{
		Title:   generateFeedTitle(movie),
		Id:      generateID(movie),
		Content: generateMovieContent(movie, tmpl),
		Link:    &feeds.Link{Href: getFirstLink(movie)},
	}

	return item
}

func generateFeedTitle(movie *Movie) string {
	if movie.Year > 0 {
		return fmt.Sprintf("%s (%d)", movie.Title, movie.Year)
	}

	return movie.Title
}

func generateMovieContent(movie *Movie, tmpl *template.Template) string {
	content := &bytes.Buffer{}
	err := tmpl.Execute(content, movie)
	if err != nil {
		fmt.Printf("failed to execute template for movie: %v\n %v", movie, err)
		return ""
	}

	return content.String()
}

func generateID(movie *Movie) string {
	titleSlug := strings.Replace(movie.Title, " ", "-", -1)
	if movie.Year > 0 {
		titleSlug = fmt.Sprintf("%s-%d", titleSlug, movie.Year)
	}

	return titleSlug
}

func getFirstLink(movie *Movie) string {
	if len(movie.PTMedias) > 0 {
		return movie.PTMedias[0].Link
	}

	return ""
}

//New creates a new Store
func New() *Store {
	return &Store{
		movies:       []Movie{},
		indexByTitle: make(map[string]*Movie),
	}
}
