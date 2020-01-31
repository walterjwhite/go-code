package main

import (
	"github.com/blevesearch/bleve"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"

	"flag"
)

type testData struct {
	Id   string
	From string
	Body string
}

var (
	indexFlag = flag.String("IndexName", "TestIndex.bleve", "Index Name to use")
)

func init() {
	application.Configure()
}

// TODO: integrate win10 / dbus notifications
func main() {
	indexData()

	search("e01")
	search("Walter")
	search("White")
	search("data")
	search("Walter.White")
}

func indexData() {
	log.Info().Msgf("Indexing data")

	d := &testData{
		Id:   "e01",
		From: "Walter.White",
		Body: "data goes here",
	}

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(*indexFlag, mapping)
	logging.Panic(err)

	index.Index(d.Id, d)

	log.Info().Msgf("Indexed data")
	logging.Panic(index.Close())
	log.Info().Msgf("Closed index")
}

func search(query string) {
	log.Info().Msgf("Search %v", query)

	index, _ := bleve.Open(*indexFlag)
	//defer logging.Panic(index.Close())

	q := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(q)
	searchResult, _ := index.Search(searchRequest)

	log.Info().Msgf("Search %v -> %v", q, searchResult)

	logging.Panic(index.Close())
}
