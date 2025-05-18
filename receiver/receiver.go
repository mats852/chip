package receiver

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/mats852/chip"
)

type ReceiverRepository interface {
	Store(ctx context.Context, chipSnapshot chip.Snapshot) error
}

type Receiver struct {
	receiverRepository ReceiverRepository
}

func NewReceiver(receiverRepository ReceiverRepository) *Receiver {
	return &Receiver{
		receiverRepository: receiverRepository,
	}
}

func (r *Receiver) Handle(data []byte) error {
	snapshot := chip.Snapshot{}

	if err := json.Unmarshal(data, &snapshot); err != nil {
		slog.Error("failed to unmarshal data", "error", err)
		return err
	}

	return r.receiverRepository.Store(context.TODO(), snapshot)
}
