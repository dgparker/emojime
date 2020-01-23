package memdb

import "github.com/hashicorp/go-memdb"

// schema is the default schema for memdb
var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"emoji": &memdb.TableSchema{
			Name: "emoji",
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:   "id",
					Unique: true,
					Indexer: &memdb.IntFieldIndex{
						Field: "No",
					},
				},
				"emoji": &memdb.IndexSchema{
					Name:   "emoji",
					Unique: true,
					Indexer: &memdb.StringFieldIndex{
						Field: "Emoji",
					},
				},
				"category": &memdb.IndexSchema{
					Name:   "category",
					Unique: false,
					Indexer: &memdb.StringFieldIndex{
						Field: "Category",
					},
				},
				"subcategory": &memdb.IndexSchema{
					Name:   "subcategory",
					Unique: false,
					Indexer: &memdb.StringFieldIndex{
						Field: "SubCategory",
					},
				},
				"unicode": &memdb.IndexSchema{
					Name:   "unicode",
					Unique: true,
					Indexer: &memdb.StringFieldIndex{
						Field: "Unicode",
					},
				},
				"name": &memdb.IndexSchema{
					Name:   "name",
					Unique: true,
					Indexer: &memdb.StringFieldIndex{
						Field: "Name",
					},
				},
				"tags": &memdb.IndexSchema{
					Name:   "tags",
					Unique: false,
					Indexer: &memdb.StringFieldIndex{
						Field: "Tags",
					},
				},
				"shortcode": &memdb.IndexSchema{
					Name:         "shortcode",
					Unique:       true,
					AllowMissing: true,
					Indexer: &memdb.StringFieldIndex{
						Field: "Shortcode",
					},
				},
			},
		},
	},
}
