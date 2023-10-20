package bulk

import (
	"github.com/olivere/elastic/v7"
	"github.com/walterjwhite/go-code/lib/data/document"
)

func (b *MasterBatch) Update(document document.Document) {
	documentTypeName, indexName := b.nodeConfiguration.PrepareIndex(document)

	request := elastic.NewBulkUpdateRequest().Index(indexName).Type(documentTypeName).Id(document.DocumentId()).Doc(document)

	b.bulkProcessor.Add(request)
}
