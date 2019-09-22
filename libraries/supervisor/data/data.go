package data

const Good = 0
const Warning = 1
const Bad = 2

type Header struct {
	Names []string
}

type Row struct {
	Status int
	Values []string
}

type DataSet struct {
	Header *Header
	Rows   []Row
}

/*
type Text struct {
	Label string
	Lines []string
}
*/
