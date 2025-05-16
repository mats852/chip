package export

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/mats852/chip"
	"github.com/mats852/chip/pkg/export/dto"
)

func Receive(data []byte) {
	responseDto := dto.GetRootAsChipDto(data, 0)

	timestamp := time.Unix(int64(responseDto.Timestamp()), 0)

	slog.Info("Receiving data", "timestamp", timestamp.Format(time.RFC3339), "raw", responseDto.Timestamp())

	for i := range responseDto.ChipsLength() {
		chipDto := &dto.Chip{}
		responseDto.Chips(chipDto, i)

		chip := &chip.Chip{
			ID: uuid.UUID(chipDto.UuidBytes()),
		}

		chip.Flags.Store(chipDto.Flags())

		slog.Info("received chip", "id", chip.ID, "flags", fmt.Sprintf("0b%064b", chip.Get()))
	}
}
