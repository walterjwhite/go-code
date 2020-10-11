package document

type Document interface {
	DocumentId() string
	//Mapping() string
}

func Equals(a Document, b Document) bool {
	// TODO: also check types
	return a.DocumentId() == b.DocumentId()
}
