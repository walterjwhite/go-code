package bulk

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/data/elasticsearch"
	"time"
)

type Operation int

const (
	Index Operation = iota
	Update
	Delete
)

type MasterBatch struct {
	nodeConfiguration *elasticsearch.NodeConfiguration
	bulkProcessor     *elastic.BulkProcessor
}

/*
type RecordOperation struct {
	Operation Operation
	DocumentType string
	DocumentId string
}

type Batch struct {
	RecordOperations []RecordOperation
}
*/

func NewDefaultBatch(c *elasticsearch.NodeConfiguration) *MasterBatch {
	return NewBatch(c, 10, 10, 1*time.Second, 2)
}

func NewBatch(c *elasticsearch.NodeConfiguration, actionSize int, dataSize int, interval time.Duration, workers int) *MasterBatch {
	masterBatch := MasterBatch{nodeConfiguration: c}

	p, err := c.Client.BulkProcessor().
		Name(getProcessorName(c)).
		Workers(workers).
		BulkActions(actionSize).
		FlushInterval(interval).
		Do(context.Background())

	logging.Panic(err)

	masterBatch.bulkProcessor = p

	return &masterBatch
}

func getProcessorName(c *elasticsearch.NodeConfiguration) string {
	return fmt.Sprintf("bulkProcessor.%v", c.IndexPrefix)
}

/*
func (b *MasterBatch) getDocumentTypeName(document interface{}) string {
	documentTypeName := typename.Get(document)
	b.nodeConfiguration.PrepareIndex(documentTypeName)

	return documentTypeName
}

func (b *MasterBatch) getIndexName(documentTypeName string) string {
	return b.nodeConfiguration.getIndexName(documentTypeName)
}
*/

func (b *MasterBatch) Flush() {
	logging.Panic(b.bulkProcessor.Flush())
}

/*
type BulkCommandFailed struct {
	RemainingCommands int
}

func (b *BulkCommandFailed) Error() string {
	return fmt.Sprintf("Bulk command has %v commands remaining, but should be 0\n", b.RemainingCommands)
}
*/
