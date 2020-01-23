package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgparker/emojime/pkg/fetcher/emojo"
)

const (
	// SAVELOCATION emojime source data save location
	SAVELOCATION = "EMOJIME_SAVE_LOCATION"
	// DEFAULTSAVELOCATION emojime default save location
	DEFAULTSAVELOCATION = "emojime.json"
)

func main() {
	saveloc := os.Getenv(SAVELOCATION)
	if saveloc == "" {
		saveloc = DEFAULTSAVELOCATION
	}

	emojoClient := emojo.NewWithDefault()

	emojis, err := emojoClient.Fetch()
	if err != nil {
		log.Fatal(err)
	}

	bdata, err := json.Marshal(emojis)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(saveloc, bdata, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
