package main

import (
	"errors"
	"log"
	"os"

	"github.com/dgparker/emojime/pkg/fetcher/static"
	"github.com/dgparker/emojime/pkg/repository/bleve"
)

const (
	// SOURCEFILE emojime source file containing all emoji data
	SOURCEFILE = "EMOJIME_SRC"
	// SEARCHINDEX emojime bleve search index containing indexed search data
	SEARCHINDEX = "EMOJIME_SEARCH_INDEX"
	// DEFAULTSEARCHINDEX emojime default bleve search index
	DEFAULTSEARCHINDEX = "emojime.bleve"
)

var (
	// ErrSrcNotSpecified is returned when EMOJIME_SRC env var is not set
	ErrSrcNotSpecified = errors.New("emojime source file not specified, please set env EMOJIME_SRC")
)

func main() {
	src := os.Getenv(SOURCEFILE)
	if src == "" {
		log.Fatal(ErrSrcNotSpecified)
	}
	index := os.Getenv(SEARCHINDEX)
	if index == "" {
		index = DEFAULTSEARCHINDEX
	}

	staticFetcher := static.New(src)

	_, err := bleve.New(index, staticFetcher)
	if err != nil {
		log.Fatal(err)
	}
}
