package elasticsearch

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/olivere/elastic/v7"
	"context"
	"time"
)

type MasterBatch struct {
	nodeConfiguration *NodeConfiguration
	indexPrepared     bool

	bulkProcessor *elastic.BulkProcessor
}

func (c *NodeConfiguration) NewBatch(actionSize int, dataSize int, interval time.Duration, workers int) *MasterBatch {
	masterBatch := MasterBatch{indexPrepared: false, nodeConfiguration: c}

	// TODO: make the size configurable
	p, err := c.Client.BulkProcessor().
		Name(c.getProcessorName()).
		Workers(workers).
		BulkActions(actionSize).                    // commit if # requests >= 1000
		BulkSize(2 << 20).                          // commit if size of requests >= 2 MB
		FlushInterval(interval /*30*time.Second*/). // commit every 30s
		Do(context.Background())

	logging.Panic(err)

	masterBatch.bulkProcessor = p

	return &masterBatch
}

func (c *NodeConfiguration) getProcessorName() string {
	return fmt.Sprintf("bulkProcessor.%v", c.IndexPrefix)
}

func (b *MasterBatch) Append(id string, document interface{}) {
	b.prepareIndex(document)

	documentTypeName := b.nodeConfiguration.getDocumentTypeName(document)

	// TODO: generalize this
	request := elastic.NewBulkIndexRequest().Index(b.nodeConfiguration.getIndexName(documentTypeName)).Type(documentTypeName).Id(id).Doc(document)
	b.bulkProcessor.Add(request)
}

func (b *MasterBatch) prepareIndex(document interface{}) {
	if !b.indexPrepared {
		b.indexPrepared = true

		// ensure this is done before the processor ever attempts, synchronous to avoid a potential race condition
		b.nodeConfiguration.prepareIndex(b.nodeConfiguration.getDocumentTypeName(document))
	}
}

func (b *MasterBatch) Flush() {
	logging.Panic(b.bulkProcessor.Flush())
}

type BulkCommandFailed struct {
	RemainingCommands int
}

func (b *BulkCommandFailed) Error() string {
	return fmt.Sprintf("Bulk command has %v commands remaining, but should be 0\n", b.RemainingCommands)
}
