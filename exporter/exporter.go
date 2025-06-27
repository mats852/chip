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

	DefaultInterval = 5 * time.Minute
)

type ExporterOpts struct {
	Interval time.Duration
}

// Exporter is responsible for exporting chips to a sender.
// It periodically sends the chips to the sender and collects the values by
// calling Clear on each chip to retrieve and reset the flags.
type Exporter struct {
	sender sender.Sender
	ticker *time.Ticker
	chips  []*chip.Chip
}

func NewExporter(sndr sender.Sender, opts ExporterOpts) (*Exporter, error) {
	return &Exporter{
		ticker: time.NewTicker(opts.Interval), // TODO: validation and default
		sender: sndr,
		chips:  nil,
	}, nil
}

func (e *Exporter) Add(c *chip.Chip, chips ...*chip.Chip) error {
	count := 1 + len(chips)

	if count+len(e.chips) > MaxChips {
		return fmt.Errorf("expects 1 to %d chips, has %d, adding %d", MaxChips, len(e.chips), count)
	}

	e.chips = append(e.chips, c)
	e.chips = append(e.chips, chips...)

	return nil
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
