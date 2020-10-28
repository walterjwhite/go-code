package elasticsearch

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/data/document"
	"github.com/walterjwhite/go/lib/utils/typename"
	"strings"
)

func (c *NodeConfiguration) PrepareIndex(document document.Document) (string, string) {
	documentTypeName, indexName := c.getIndexName(document)

	_, exists := c.Indexes[documentTypeName]
	if !exists {
		log.Info().Msgf("Index: %v/%v\n", documentTypeName, indexName)

		c.doPrepareIndex(indexName)
		c.Indexes[documentTypeName] = true
	}

	return documentTypeName, indexName
}

func (c *NodeConfiguration) doPrepareIndex(indexName string) {
	if c.isIndexExisting(indexName) {
		if c.DropIndex {
			c.deleteIndex(indexName)
		}
	} else {
		c.createIndex(indexName)
	}
}

func (c *NodeConfiguration) isIndexExisting(indexName string) bool {
	exists, err := c.Client.IndexExists(indexName).Do(context.Background())
	logging.Panic(err)

	return exists
}

func (c *NodeConfiguration) createIndex(indexName string) {
	createIndexService := c.Client.CreateIndex(indexName)

	mapping, exists := c.Mappings[indexName]
	if exists {
		createIndexService.BodyString(mapping)
	}

	/*result*/
	_, err := createIndexService.Do(context.Background())
	logging.Panic(err)
}

func (c *NodeConfiguration) deleteIndex(indexName string) {
	/*result*/ _, err := c.Client.DeleteIndex(indexName).Do(context.Background())
	logging.Panic(err)
}

func (c *NodeConfiguration) getIndexName(document document.Document) (string, string) {
	documentTypeName := strings.ToLower(typename.Get(document))
	return documentTypeName, c.getFullIndexName(documentTypeName)
}

func (c *NodeConfiguration) getFullIndexName(documentTypeName string) string {
	if len(c.IndexPrefix) > 0 {
		return strings.ToLower(fmt.Sprintf("%v.%v", c.IndexPrefix, documentTypeName))
	}

	return strings.ToLower(documentTypeName)
}
