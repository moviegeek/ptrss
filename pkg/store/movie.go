package store

import "github.com/moviegeek/pt"

//Movie is the movie info, with all downloading source from differnet pt sites
type Movie struct {
	Title     string
	Year      int
	IMDBID    string
	Published string
	Updated   string
	PTMedias  []PTMedia
}

//PTMedia media metadata of a pt movie
type PTMedia struct {
	pt.MovieInfo
	Site       string
	Link       string
	TorrentURL string
	SiteID     string
}
