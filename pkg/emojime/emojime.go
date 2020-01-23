package emojime

// Fetcher defines the interface to load emoji source data
type Fetcher interface {
	Fetch() ([]*Emoji, error)
}

// Reloader defines the interface to reload emoji source data
type Reloader interface {
	Reload() error
}

// Searcher defines the interface to reload and search the emoji source data
type Searcher interface {
	Reloader
	Search(query string) ([]string, error)
}

// Repository defines the interface to reload and interact with emoji source data
type Repository interface {
	Reloader
	List() ([]*Emoji, error)
	Get(name string) (*Emoji, error)
}

// Emoji represents the data structure for an emoji character
type Emoji struct {
	No          int    `json:"No"`
	Emoji       string `json:"Emoji"`
	Category    string `json:"Category"`
	SubCategory string `json:"SubCategory"`
	Unicode     string `json:"Unicode"`
	Name        string `json:"Name"`
	Tags        string `json:"Tags"`
	Shortcode   string `json:"Shortcode"`
}

// Client represents an emoji client and acts as a layer of abstraction between the searcher
// and repostiory interfaces in an effort to maintain a stable api
type Client struct {
	repo  Repository
	index Searcher
}

// New takes a Repository implementation and returns an emoji Client
func New(r Repository, i Searcher) *Client {
	return &Client{
		repo:  r,
		index: i,
	}
}

// ListEmojis fetches a list of emojis from the repository
func (c *Client) ListEmojis() ([]*Emoji, error) {
	return c.repo.List()
}

// GetEmoji fetches an emoji from the repository where name matches
func (c *Client) GetEmoji(name string) (*Emoji, error) {
	return c.repo.Get(name)
}

// SearchEmojis ...
func (c *Client) SearchEmojis(query string) ([]*Emoji, error) {
	matchIDs, err := c.index.Search(query)
	if err != nil {
		return nil, err
	}

	var emojis []*Emoji
	for _, v := range matchIDs {
		emoji, err := c.repo.Get(v)
		if err != nil {
			return nil, err
		}

		emojis = append(emojis, emoji)
	}
	return emojis, nil
}

// UpdateEmojis TO BE IMPLEMENTED
func (c *Client) UpdateEmojis() error {
	var err error

	err = c.repo.Reload()
	if err != nil {
		return err
	}

	err = c.index.Reload()
	if err != nil {
		return err
	}

	return nil
}
