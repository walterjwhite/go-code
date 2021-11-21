package elasticsearch

import (
	"github.com/olivere/elastic/v7"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

/*
type Document interface {
	DocumentId() string
}
*/

type NodeConfiguration struct {
	Client *elastic.Client

	IndexPrefix string
	DropIndex   bool
	Mappings    map[string]string
	Indexes     map[string]bool
}

// unnecessary, this is the default
func NewDefaultClient() *NodeConfiguration {
	nodeConfiguration := NodeConfiguration{}
	nodeConfiguration.configure()

	return &nodeConfiguration
}

func (c *NodeConfiguration) configure() {
	// Create a client and connect to http://127.0.0.1:9200 (default)
	client, err := elastic.NewClient()
	logging.Panic(err)

	c.Client = client

	if c.Indexes == nil {
		c.Indexes = make(map[string]bool)
	}
}
