package export

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/mats852/chip"
	"github.com/mats852/chip/pkg/export/dto"
)

type Exporter struct {
	client *http.Client
	ticker *time.Ticker

	chips []*chip.Chip
}

func NewExporter(c ...*chip.Chip) *Exporter {
	return &Exporter{
		client: &http.Client{},
		ticker: time.NewTicker(5 * time.Second),
		chips:  c,
	}
}

// Serve starts the exporter and periodically sends the chips.
func (e *Exporter) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-e.ticker.C:
			e.export()
		}
	}
}

func (e *Exporter) export() {
	slog.Info("exporting chips")

	builder := flatbuffers.NewBuilder(64 + len(e.chips)*32) // re-evaluate this size

	chipOffsets := make([]flatbuffers.UOffsetT, len(e.chips))

	for i, chip := range e.chips {
		idPos := builder.CreateByteVector(chip.ID[:])

		flag := chip.Clear()

		dto.ChipStart(builder)
		dto.ChipAddUuid(builder, idPos)
		dto.ChipAddFlags(builder, flag)

		chipOffsets[i] = dto.ChipEnd(builder)

		slog.Info("exported chip", "id", chip.ID, "flags", fmt.Sprintf("0b%064b", flag))
	}

	dto.ChipDtoStartChipsVector(builder, len(chipOffsets))
	for i := len(chipOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(chipOffsets[i])
	}

	chipsVector := builder.EndVector(len(chipOffsets))

	dto.ChipDtoStart(builder)
	dto.ChipDtoAddTimestamp(builder, uint64(time.Now().Unix()))
	dto.ChipDtoAddChips(builder, chipsVector)

	builder.Finish(dto.ChipDtoEnd(builder))

	// TODO: send on the wire to the listener
	Receive(builder.FinishedBytes())
}
