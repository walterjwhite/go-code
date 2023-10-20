package bulk

import (
	"github.com/olivere/elastic/v7"
	"github.com/walterjwhite/go-code/lib/data/document"
)

func (b *MasterBatch) Delete(document document.Document) {
	documentTypeName, indexName := b.nodeConfiguration.PrepareIndex(document)

	request := elastic.NewBulkDeleteRequest().Index(indexName).Type(documentTypeName).Id(document.DocumentId())

	b.bulkProcessor.Add(request)
}
