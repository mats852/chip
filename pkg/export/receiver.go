package export

import (
	"context"
  "log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/mats852/chip"
	"github.com/mats852/chip/pkg/export/dto"
)

type ReceiverRepository interface {
	Store(ctx context.Context, timestamp time.Time, chip *chip.Chip) error
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
	responseDto := dto.GetRootAsChipDto(data, 0)

	timestamp := time.Unix(int64(responseDto.Timestamp()), 0)

	slog.Info("receiving data", "time", timestamp.Format(time.RFC3339), "timestamp", responseDto.Timestamp())

	for i := range responseDto.ChipsLength() {
		chipDto := &dto.Chip{}
		responseDto.Chips(chipDto, i)

		chip := &chip.Chip{
			ID: uuid.UUID(chipDto.UuidBytes()),
		}

		chip.Flags.Store(chipDto.Flags())

    r.receiverRepository.Store(context.TODO(), timestamp, chip)
	}

  return nil
}
