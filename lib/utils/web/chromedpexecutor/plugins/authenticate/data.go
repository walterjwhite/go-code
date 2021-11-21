package authenticate

import (
	"time"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
)

type Session struct {
	Credentials *Credentials
	Website     *Website
	//IsKeepAlive bool

	chromedpsession *chromedpexecutor.ChromeDPSession
}

type Credentials struct {
	Secrets []*FieldSecret
}

type FieldSecret struct {
	// must match Field.id to be injected
	FieldId     *string
	SecretKey   *string
	SecretValue *string
}

type Website struct {
	Url            *string
	FieldGroups    []*FieldGroup
	SessionTimeout *time.Duration
	KeepAliveUrl   *string
}

type FieldGroup struct {
	Fields         []*Field
	SubmitSelector *string
}

type Field struct {
	Id *string

	Selector *string
}

// TODO: add support for entering one-time tokens
// TODO: add support for answering challenge questions
