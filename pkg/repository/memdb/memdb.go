package memdb

import (
	"github.com/dgparker/emojime/pkg/emojime"
	"github.com/hashicorp/go-memdb"
)

const (
	emojiTable = "emoji"
)

// Client implements emoji.Repository interface by utlizizing go-memdb
type Client struct {
	db *memdb.MemDB
}

// New returns a new DB struct for accessing the memdb repository
// src should be a file path to the source data file to load into memory
func New(fetcher emojime.Fetcher) (*Client, error) {
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	c := &Client{
		db: db,
	}

	err = c.load(fetcher)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// load initializes memdb using the provided fetcher
func (c *Client) load(fetcher emojime.Fetcher) error {
	emojis, err := fetcher.Fetch()
	if err != nil {
		return err
	}

	txn := c.db.Txn(true)
	defer txn.Abort()

	for _, e := range emojis {
		if err := txn.Insert(emojiTable, e); err != nil {
			return err
		}
	}

	txn.Commit()

	return nil
}

// List returns all emojis in the repository
func (c *Client) List() ([]*emojime.Emoji, error) {
	txn := c.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get(emojiTable, "id")
	if err != nil {
		return nil, err
	}

	var result []*emojime.Emoji
	for obj := it.Next(); obj != nil; obj = it.Next() {
		result = append(result, obj.(*emojime.Emoji))
	}

	return result, nil
}

// Get returns an emoji where name matches repository index name
func (c *Client) Get(name string) (*emojime.Emoji, error) {
	txn := c.db.Txn(false)
	defer txn.Abort()

	result, err := txn.First(emojiTable, "name", name)
	if err != nil {
		return nil, err
	}

	return result.(*emojime.Emoji), nil
}

// Reload TO BE IMPLEMENTED
func (c *Client) Reload() error {
	return nil
}
