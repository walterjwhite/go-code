package notification

import (
	"github.com/walterjwhite/go/lib/application/property"
	"github.com/walterjwhite/go/lib/io/yaml"
	"os"
	"path/filepath"
)

type ReferenceData struct {
	Filenames []string
}

func (n *Notification) loadReferences() {
	referencesFilename := n.getReferencesFilename()

	if !exists(referencesFilename) {
		return
	}

	r := &ReferenceData{}
	yaml.Read(referencesFilename, r)

	if len(r.Filenames) > 0 {
		n.Filenames = append(n.Filenames, r.Filenames...)
	}
}

func (n *Notification) getReferencesFilename() string {
	return filepath.Join(property.GetConfigurationDirectory("spot", n.Session.FeedId, "notifications"), "references.yaml")
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
