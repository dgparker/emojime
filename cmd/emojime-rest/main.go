package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dgparker/emojime/pkg/emojime"
	"github.com/dgparker/emojime/pkg/fetcher/static"
	"github.com/dgparker/emojime/pkg/http/rest"
	"github.com/dgparker/emojime/pkg/repository/bleve"
	"github.com/dgparker/emojime/pkg/repository/memdb"
	"go.uber.org/zap"
)

const (
	// PORT emojime port env name
	PORT = "EMOJIME_PORT"
	// DEFAULTPORT emojime default port value
	DEFAULTPORT = "9000"
	// ENV emojime env value
	ENV = "EMOJIME_ENV"
	// DEFAULTENV emojime default env value
	DEFAULTENV = "dev"
	// SOURCEFILE emojime source file location containg all emoji data
	SOURCEFILE = "EMOJIME_SRC"
	// SEARCHINDEX emojime bleve search index containing indexed search data
	SEARCHINDEX = "EMOJIME_SEARCH_INDEX"
	// DEFAULTSEARCHINDEX emojime default bleve search index
	DEFAULTSEARCHINDEX = "emojime.bleve"
)

func main() {
	port := os.Getenv(PORT)
	if port == "" {
		port = DEFAULTPORT
	}

	env := os.Getenv(ENV)
	if env == "" {
		env = DEFAULTENV
	}

	src := os.Getenv(SOURCEFILE)
	if src == "" {
		log.Fatal("source file required for initialization please set $EMOJIME_SRC")
	}

	index := os.Getenv(SEARCHINDEX)
	if index == "" {
		index = DEFAULTSEARCHINDEX
	}

	var logger *zap.Logger
	var err error
	switch env {
	case "dev":
		logger, err = zap.NewDevelopment()
	case "prod":
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal(err)
	}

	staticFetcher := static.New(src)

	repo, err := memdb.New(staticFetcher)
	if err != nil {
		log.Fatal(err)
	}

	search, err := bleve.New(index, staticFetcher)
	if err != nil {
		log.Fatal(err)
	}

	svc := emojime.New(repo, search)

	api := rest.New(logger, svc)
	api.Addr = fmt.Sprintf(":%s", port)
	logger.Info(
		"server start",
		zap.String("port", port),
	)
	log.Fatal(api.ListenAndServe())
}
