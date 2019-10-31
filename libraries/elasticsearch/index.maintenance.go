package elasticsearch

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
	"log"
	"strings"
)

func (c *NodeConfiguration) getDocumentTypeName(document interface{}) string {
	return strings.ToLower(strings.ReplaceAll(fmt.Sprintf("%T", document), "*", ""))
}

func (c *NodeConfiguration) prepareIndex(documentTypeName string) {
	indexName := c.getIndexName(documentTypeName)

	log.Printf("Index: %v/%v\n", documentTypeName, indexName)

	if c.isIndexExisting(indexName) {
		if !c.DropIndex {
			return
		}

		c.deleteIndex(indexName)
	}

	c.createIndex(indexName)
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

func (c *NodeConfiguration) getIndexName(documentTypeName string) string {
	if len(c.IndexPrefix) > 0 /* != nil*/ {
		return strings.ToLower(fmt.Sprintf("%v.%v", c.IndexPrefix, documentTypeName))
	}

	return strings.ToLower(documentTypeName)
}
