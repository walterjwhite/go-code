package bulk

import (
	"github.com/olivere/elastic/v7"
	"github.com/walterjwhite/go/lib/data/document"
)

func (b *MasterBatch) Index(document document.Document) {
	documentTypeName, indexName := b.nodeConfiguration.PrepareIndex(document)

	request := elastic.NewBulkIndexRequest().Index(indexName).Type(documentTypeName).Id(document.DocumentId()).Doc(document)

	b.bulkProcessor.Add(request)
}
