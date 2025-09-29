package citrix

import (
	"context"
)

type PromptHandler interface {
	Handle(ctx context.Context)
}
