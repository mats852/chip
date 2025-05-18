package sender

import (
	"context"

	"github.com/mats852/chip"
)

type Sender interface {
	Send(context.Context, chip.Snapshot) error
}
