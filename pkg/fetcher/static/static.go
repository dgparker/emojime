package static

import (
	"encoding/json"
	"github.com/dgparker/emojime/pkg/emojime"
	"io/ioutil"
)

// Client represents a static client which implements the emojime.Fetcher interface
// mainly used for local testing/debugging which will fetch emoji data from a local file
type Client struct {
	SourcePath string
}

// New returns a new static fetcher client
func New(sourcePath string) *Client {
	return &Client{
		SourcePath: sourcePath,
	}
}

// Fetch retrieves emoji data from the source path specified in the client configuration
func (c *Client) Fetch() ([]*emojime.Emoji, error) {
	data, err := ioutil.ReadFile(c.SourcePath)
	if err != nil {
		return nil, err
	}

	var emojis []*emojime.Emoji
	err = json.Unmarshal(data, &emojis)
	if err != nil {
		return nil, err
	}

	return emojis, nil
}
