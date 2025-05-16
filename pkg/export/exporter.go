package export

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/mats852/chip"
	"github.com/mats852/chip/pkg/export/dto"
)

const (
	MaxChips = 64 // arbitrary number for now
)

type ExporterSender interface {
	Send(ctx context.Context, data []byte) error
}

// Exporter is responsible for exporting chips to a sender.
// It periodically sends the chips to the sender and collects the values by
// calling Clear on each chip to retrieve and reset the flags.
type Exporter struct {
	sender ExporterSender
	ticker *time.Ticker
	chips  []*chip.Chip
}

func NewExporter(sender ExporterSender, c ...*chip.Chip) (*Exporter, error) {
	if len(c) == 0 || len(c) > MaxChips {
		return nil, fmt.Errorf("expects 1 to %d chips, received %d", MaxChips, len(c))
	}

	return &Exporter{
		ticker: time.NewTicker(5 * time.Second),
		sender: sender,
		chips:  c,
	}, nil
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
	slog.Debug("exporting chips")

	data := serializeChips(e.chips)

	// TODO: 1. enqueue in buffer and dequeue the buffer with the client
	// TODO: 2. send on the wire to the listener
	// TODO: 3. add an in memory buffer to keep n records in case the network is down

	if err := e.sender.Send(context.TODO(), data); err != nil {
		slog.Error("failed to send data", "error", err)
	}
}

func serializeChips(chips []*chip.Chip) []byte {
	builder := flatbuffers.NewBuilder(64 + len(chips)*32) // re-evaluate this size

	chipOffsets := make([]flatbuffers.UOffsetT, len(chips))

	for i, chip := range chips {
		idPos := builder.CreateByteVector(chip.ID[:])

		dto.ChipStart(builder)
		dto.ChipAddUuid(builder, idPos)
		dto.ChipAddFlags(builder, chip.Clear())

		chipOffsets[i] = dto.ChipEnd(builder)
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

	return builder.FinishedBytes()
}
