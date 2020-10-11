package sourcecode

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"github.com/walterjwhite/go-application/libraries/data/elasticsearch"
	"github.com/walterjwhite/go-application/libraries/data/elasticsearch/bulk"
	"github.com/walterjwhite/go-application/libraries/utils/foreachfile"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type IndexSession struct {
	Path string
	Tags []string

	masterBatch *bulk.MasterBatch
	waitGroup   *sync.WaitGroup

	waitChannel chan bool
}

// TODO:
// 1. delete old contents under path, tags first
// 2. create index with mapping
func (i *IndexSession) Prune() {
	// find all matching documents
	// delete
}

func (i *IndexSession) IndexRecursive() {
	i.masterBatch = bulk.NewDefaultBatch(elasticsearch.NewDefaultClient())
	defer i.masterBatch.Flush()

	i.waitGroup = &sync.WaitGroup{}
	i.waitChannel = make(chan bool, 8)

	foreachfile.ExecuteCallback(i.Path, i.Index, &foreachfile.HiddenFileExcluder{})

	i.waitGroup.Wait()
	//<- i.waitChannel
}

func (i *IndexSession) Index(filePath string) {
	i.waitGroup.Add(1)
	go i.indexFile(filePath)
}

func (i *IndexSession) indexFile(filePath string) {
	i.waitChannel <- true
	defer i.waitGroup.Done()

	path := filepath.Dir(filePath)
	name := filepath.Base(filePath)

	// TODO: make this configurable
	if strings.Index(name, ".") == 0 {
		<-i.waitChannel
		log.Info().Msgf("skipping: %v", name)
		return
	}

	log.Info().Msgf("indexing: %v", name)

	file, err := os.Open(filePath)
	logging.Panic(err)

	defer file.Close()

	number := 0

	waitGroup := &sync.WaitGroup{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := &Line{Path: path, Name: name, Number: number, Tags: i.Tags, Contents: scanner.Text()}
		number++

		waitGroup.Add(1)
		go i.doIndex(waitGroup, l)
	}

	err = scanner.Err()
	logging.Panic(err)

	waitGroup.Wait()
	<-i.waitChannel
}

func (i *IndexSession) doIndex(waitGroup *sync.WaitGroup, l *Line) {
	defer waitGroup.Done()

	i.masterBatch.Index(l)
}
