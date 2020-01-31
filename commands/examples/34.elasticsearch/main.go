package main

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/elasticsearch"
	"strconv"
	"time"
)

type TestData struct {
	Id          string
	Name        string
	Description string
}

func init() {
	application.Configure()
}

func main() {
	c := elasticsearch.NewDefaultClient()
	b := c.NewBatch(10, 5242880, 30*time.Second, 2)

	for i := 1; i <= 100; i++ {
		d := TestData{Id: strconv.Itoa(i + 100), Name: fmt.Sprintf("2019/12/18 %v", i), Description: fmt.Sprintf("D %v", i)}

		b.Index(d.Id, d)
	}

	b.Flush()

	/*
		for i := 1; i <= 100; i++ {
			d := TestData{Id: strconv.Itoa(i), Name: fmt.Sprintf("N %v", i), Description: fmt.Sprintf("D %v", i)}

			b.Delete(d.Id, d)
		}

		b.Flush()
	*/
}
