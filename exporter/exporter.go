package exporter

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/mats852/chip"
	"github.com/mats852/chip/sender"
)

const (
	MaxChips = 64 // arbitrary number for now
)

// Exporter is responsible for exporting chips to a sender.
// It periodically sends the chips to the sender and collects the values by
// calling Clear on each chip to retrieve and reset the flags.
type Exporter struct {
	sender sender.Sender
	ticker *time.Ticker
	chips  []*chip.Chip
}

func NewExporter(sndr sender.Sender, c ...*chip.Chip) (*Exporter, error) {
	if len(c) == 0 || len(c) > MaxChips {
		return nil, fmt.Errorf("expects 1 to %d chips, received %d", MaxChips, len(c))
	}

	return &Exporter{
		ticker: time.NewTicker(5 * time.Second),
		sender: sndr,
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
	if err := e.sender.Send(context.TODO(), chip.NewSnapshot(e.chips)); err != nil {
		slog.Error("failed to send data", "error", err)
	}
}
