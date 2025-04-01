package sourcecode

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/data/elasticsearch"
	"github.com/walterjwhite/go-code/lib/data/elasticsearch/bulk"
	"github.com/walterjwhite/go-code/lib/utils/foreachfile"
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

func (i *IndexSession) Prune() {
}

func (i *IndexSession) IndexRecursive() {
	i.masterBatch = bulk.NewDefaultBatch(elasticsearch.NewDefaultClient())
	defer i.masterBatch.Flush()

	i.waitGroup = &sync.WaitGroup{}
	i.waitChannel = make(chan bool, 8)

	foreachfile.ExecuteCallback(i.Path, i.Index, &foreachfile.HiddenFileExcluder{})

	i.waitGroup.Wait()
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

	if strings.Index(name, ".") == 0 {
		<-i.waitChannel
		log.Info().Msgf("skipping: %v", name)
		return
	}

	log.Info().Msgf("indexing: %v", name)

	file, err := os.Open(filePath)
	logging.Panic(err)

	defer close(file)

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

func close(file *os.File) {
	logging.Panic(file.Close())
}

func (i *IndexSession) doIndex(waitGroup *sync.WaitGroup, l *Line) {
	defer waitGroup.Done()

	i.masterBatch.Index(l)
}
