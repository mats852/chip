package export

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mats852/chip"
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

	for _, chip := range e.chips {
		flag := chip.Clear()
		slog.Info("exported chip", "id", chip.ID, "flags", fmt.Sprintf("0b%064b", flag))
	}

	// TODO: send on the wire to the listener
}
