package store

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/moviegeek/pt"
)

func TestParser(t *testing.T) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://hdchina.org/torrentrss.php?rows=50&cat17=1&cat9=1&isize=1&passkey=0183e8204f352824ec28e6fe7c132871")
	if err != nil {
		t.Error(err)
	}

	for _, i := range feed.Items {
		printItem(i)
	}
}

func printItem(item *gofeed.Item) {
	content := item.Content
	item.Content = fmt.Sprintf("content length: %d", len(content))
	description := item.Description
	item.Description = fmt.Sprintf("description length: %d", len(description))

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", pt.ParseTitle(item.Title))

	item.Content = content
	item.Description = description
}

func TestAddNewMovie(t *testing.T) {
	info := pt.MovieInfo{
		Title:      "The Dawn Wall",
		Year:       2017,
		Group:      "HANDJOB",
		Source:     pt.Blueray,
		Resolution: pt.HD,
		Size:       8160000000,
	}
	publishTime, _ := time.Parse(time.RFC3339, "2019-03-02T04:37:17Z0000")
	item := &gofeed.Item{
		Title:           "The.Dawn.Wall.2017.1080p.BluRay.x264-HANDJOB[8.16 GB]",
		Description:     "description length: 3129",
		Content:         "content length: 0",
		Link:            "https://hdchina.org/details.php?id=309608",
		Published:       "Sat, 02 Mar 2019 12:37:17 +0800",
		PublishedParsed: &publishTime,
		GUID:            "51e4b14ddcfdb3860a72dace87355f5878b75f5c",
		Categories:      []string{"[电影Movie(1080p)]"},
	}

	s := New()
	s.addNewMovie(info, item)

	newMovie := Movie{
		Title:     "The Dawn Wall",
		Year:      2017,
		IMDBID:    "",
		Published: "Sat, 02 Mar 2019 12:37:17 +0800",
		Updated:   "",
		PTMedias: []PTMedia{
			{
				MovieInfo:  info,
				Site:       hdcSiteName,
				Link:       "https://hdchina.org/details.php?id=309608",
				TorrentURL: "",
				SiteID:     "309608",
			},
		},
	}

	if !reflect.DeepEqual(s.movies, []Movie{newMovie}) {
		t.Fatalf("added movies does not match")
	}
}

func TestAddExistingMovie(t *testing.T) {
	s := New()
	s.movies = []Movie{
		{
			Title:     "The Dawn Wall",
			Year:      2017,
			IMDBID:    "",
			Published: "Fri, 15 Mar 2019 20:24:59 +0800",
			Updated:   "",
			PTMedias: []PTMedia{
				{
					MovieInfo: pt.MovieInfo{
						Title:      "Aquaman",
						Year:       2018,
						Group:      "HDChina",
						Source:     pt.Blueray,
						Resolution: pt.FHD,
						Size:       19600000000,
					},
					Site:       hdcSiteName,
					Link:       "https://hdchina.org/details.php?id=310862",
					TorrentURL: "",
					SiteID:     "310862",
				},
			},
		},
	}

	info := pt.MovieInfo{
		Title:      "Aquaman",
		Year:       2018,
		Group:      "HDChina",
		Source:     pt.Blueray,
		Resolution: pt.HD,
		Size:       7030000000,
	}
	publishTime, _ := time.Parse(time.RFC3339, "2019-03-15T05:49:47Z0000")
	item := &gofeed.Item{
		Title:           "Aquaman.2018.720p.BluRay.x264.DTS-HDChina[7.03 GB]",
		Description:     "description length: 7595",
		Content:         "content length: 0",
		Link:            "https://hdchina.org/details.php?id=310834",
		Published:       "Fri, 15 Mar 2019 13:49:47 +0800",
		PublishedParsed: &publishTime,
		GUID:            "a6a61d07c9de140c0aea46cb9f95e752c66a3e71",
		Categories:      []string{"[电影Movie(720p)]"},
	}

	s.addExsitMovieFromFeedItem(&s.movies[0], info, item)

	if len(s.movies) != 1 {
		t.Fatalf("expect only 1 movie in store, but got %d", len(s.movies))
	}

	newMovies := []Movie{
		{
			Title:     "The Dawn Wall",
			Year:      2017,
			IMDBID:    "",
			Published: "Fri, 15 Mar 2019 20:24:59 +0800",
			Updated:   "",
			PTMedias: []PTMedia{
				{
					MovieInfo: pt.MovieInfo{
						Title:      "Aquaman",
						Year:       2018,
						Group:      "HDChina",
						Source:     pt.Blueray,
						Resolution: pt.FHD,
						Size:       19600000000,
					},
					Site:       hdcSiteName,
					Link:       "https://hdchina.org/details.php?id=310862",
					TorrentURL: "",
					SiteID:     "310862",
				},
				{
					MovieInfo: pt.MovieInfo{
						Title:      "Aquaman",
						Year:       2018,
						Group:      "HDChina",
						Source:     pt.Blueray,
						Resolution: pt.HD,
						Size:       7030000000,
					},
					Site:       hdcSiteName,
					Link:       "https://hdchina.org/details.php?id=310834",
					TorrentURL: "",
					SiteID:     "310834",
				},
			},
		},
	}

	if !reflect.DeepEqual(s.movies, newMovies) {
		t.Fatalf("added movies does not match\n Expected: %+v\nGot: %+v", newMovies, s.movies)
	}
}

func TestToRss(t *testing.T) {
	f, err := os.OpenFile("rss.xml", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatal(err)
	}

	s := New()
	s.movies = []Movie{
		{
			Title:     "The Dawn Wall",
			Year:      2017,
			IMDBID:    "",
			Published: "Fri, 15 Mar 2019 20:24:59 +0800",
			Updated:   "",
			PTMedias: []PTMedia{
				{
					MovieInfo: pt.MovieInfo{
						Title:      "Aquaman",
						Year:       2018,
						Group:      "HDChina",
						Source:     pt.Blueray,
						Resolution: pt.FHD,
						Size:       19600000000,
					},
					Site:       hdcSiteName,
					Link:       "https://hdchina.org/details.php?id=310862",
					TorrentURL: "",
					SiteID:     "310862",
				},
			},
		},
	}

	fmt.Println(s.ToRss(f))
}
