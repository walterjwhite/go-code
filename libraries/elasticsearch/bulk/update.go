package bulk

import (
	"github.com/walterjwhite/go-application/libraries/document"
	"gopkg.in/olivere/elastic.v7"
)

func (b *MasterBatch) Update(document document.Document) {
	documentTypeName, indexName := b.nodeConfiguration.PrepareIndex(document)

	request := elastic.NewBulkUpdateRequest().Index(indexName).Type(documentTypeName).Id(document.DocumentId()).Doc(document)

	b.bulkProcessor.Add(request)
}
