package emojo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgparker/emojime/pkg/emojime"
)

const (
	// DefaultSourceURL is the default source url for the emojo emoji data
	DefaultSourceURL = "https://raw.githubusercontent.com/CodeFreezr/emojo/master/db/v5/emoji-v5.json"
)

// Client represents an emojo client which implements the emojime.Fetcher interface
type Client struct {
	SourceURL  string
	HTTPClient *http.Client
}

// New returns a new emojo client
func New(sourceURL string, client *http.Client) *Client {
	return &Client{
		SourceURL:  sourceURL,
		HTTPClient: client,
	}
}

// NewWithDefault returns a new emojo client with sane defaults
func NewWithDefault() *Client {
	c := &http.Client{
		Timeout: time.Second * 10,
	}
	return New(DefaultSourceURL, c)
}

// Fetch retrieves new emoji data from the emojo source url
func (c *Client) Fetch() ([]*emojime.Emoji, error) {
	var emojis []*emojime.Emoji

	req, err := http.NewRequest(http.MethodGet, c.SourceURL, nil)
	if err != nil {
		return emojis, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return emojis, err
	}

	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return emojis, err
	}

	err = json.Unmarshal(data, &emojis)
	if err != nil {
		return emojis, err
	}

	return emojis, nil
}
