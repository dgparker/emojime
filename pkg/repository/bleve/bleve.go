package bleve

import (
	"log"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search"
	"github.com/dgparker/emojime/pkg/emojime"
)

// Client represents a bleve client which indexes emojis to be used for search functionality
type Client struct {
	emojime.Fetcher
	index bleve.Index
}

// New returns a new bleve repository client
func New(indexPath string, fetcher emojime.Fetcher) (*Client, error) {
	i, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Println("creating new index")
		i, err = createIndex(indexPath, fetcher)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		log.Println("using existing index")
	}

	return &Client{
		Fetcher: fetcher,
		index:   i,
	}, nil
}

// Search returns a limited list of ids that match the query provided sorted by relevance
func (c *Client) Search(query string) ([]string, error) {
	q := bleve.NewQueryStringQuery(query)
	req := bleve.NewSearchRequestOptions(q, 100, 0, false)
	res, err := c.index.Search(req)
	if err != nil {
		return nil, err
	}

	matchIDs := getMatchIDs(res.Hits)

	return matchIDs, nil
}

// Reload TO BE IMPLEMENTED re-index with zero downtime
func (c *Client) Reload() error {
	return nil
}

func getMatchIDs(hits search.DocumentMatchCollection) []string {
	var matchIDs []string
	for _, v := range hits {
		matchIDs = append(matchIDs, v.ID)
	}
	return matchIDs
}

func createIndex(indexPath string, fetcher emojime.Fetcher) (bleve.Index, error) {
	indexMapping, err := buildIndexMapping()
	if err != nil {
		return nil, err
	}

	emojiIndex, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		return nil, err
	}

	err = indexEmojis(emojiIndex, fetcher)
	if err != nil {
		return nil, err
	}

	return emojiIndex, nil
}

func buildIndexMapping() (mapping.IndexMapping, error) {
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a gneric reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	emojiMapping := bleve.NewDocumentMapping()

	// category
	emojiMapping.AddFieldMappingsAt("Category", keywordFieldMapping)
	// sub-category
	emojiMapping.AddFieldMappingsAt("SubCategory", englishTextFieldMapping)
	// unicode
	emojiMapping.AddFieldMappingsAt("Unicode", englishTextFieldMapping)
	// name
	emojiMapping.AddFieldMappingsAt("Name", englishTextFieldMapping)
	// tags
	emojiMapping.AddFieldMappingsAt("Tags", keywordFieldMapping)
	// shortcode
	emojiMapping.AddFieldMappingsAt("Shortcode", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("emoji", emojiMapping)

	indexMapping.TypeField = "Name"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}

func indexEmojis(index bleve.Index, fetcher emojime.Fetcher) error {
	emojis, err := fetcher.Fetch()
	if err != nil {
		return err
	}

	batch := index.NewBatch()
	batchSize := 100
	batchCount := 0
	for _, v := range emojis {
		err := batch.Index(v.Name, v)
		if err != nil {
			return err
		}
		batchCount++

		// flush current batch
		if batchCount >= batchSize {
			err = index.Batch(batch)
			if err != nil {
				return err
			}
			batch = index.NewBatch()
			batchCount = 0
		}
	}

	// flush the last batch
	if batchCount > 0 {
		err = index.Batch(batch)
		if err != nil {
			return err
		}
	}

	log.Println("indexing complete")
	return nil
}
