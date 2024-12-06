package document

type Document interface {
	DocumentId() string
}

func Equals(a Document, b Document) bool {
	return a.DocumentId() == b.DocumentId()
}
