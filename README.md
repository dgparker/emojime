# emojime

## Install and run
```
$ git clone https://github.com/dgparker/emojime
$ cd emojime && go build -o emojime-rest cmd/emojime-rest/main.go


// run
// NOTE: default dataset and search index can be found in emojime/data folder
$ EXPORT EMOJIME_SRC=/path/to/emojime/src
$ EXPORT EMOJIME_SEARCH_INDEX=/path/to/emojime/index
$ ./emojime-rest
```

### About
emojime is a simple REST API for emoji search. 

The goal of the project is to provide a better search experience for emojis. Current implementations fail to properly group/organize emoji data in a consistent manner and often times return less than optimal results.

Emojis are currently indexed using the following fields

given the following emoji:  🎱
- Category - ```Activities```
- SubCategory - ```sport```
- Unicode character - ```U+1F3B1```
- Name - ```pool 8 ball```
- Tags - ```8, 8 ball, ball, billiard, eight, game```
- Shortcode - ```:8Ball:```

Currently using [bleve](https://github.com/blevesearch/bleve) for search and [go-memdb](https://github.com/hashicorp/go-memdb) for in memory storage. These technologies enable the emojime service to provide fast response times and consistent results.


### Future enhancements

- develop our own dataset using the official unicode emoji list found [here](https://unicode.org/emoji/charts/full-emoji-list.html)
- enable hot reload of source data
