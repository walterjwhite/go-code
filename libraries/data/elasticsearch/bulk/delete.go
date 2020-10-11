package bulk

import (
	"github.com/walterjwhite/go-application/libraries/data/document"
	"gopkg.in/olivere/elastic.v7"
)

func (b *MasterBatch) Delete(document document.Document) {
	documentTypeName, indexName := b.nodeConfiguration.PrepareIndex(document)

	request := elastic.NewBulkDeleteRequest().Index(indexName).Type(documentTypeName).Id(document.DocumentId())

	b.bulkProcessor.Add(request)
}
