package export

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mats852/chip"
)

type Exporter struct {
	client *http.Client
	ticker *time.Ticker

	chip *chip.Chip
}

func NewExporter(
	id uuid.UUID,
	c *chip.Chip,
) *Exporter {
	return &Exporter{
		client: &http.Client{},
		ticker: time.NewTicker(5 * time.Second),
		chip:   c,
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
	snapshot := e.chip.Export()

	slog.Info("exporting chips")

	for k, v := range snapshot.Records {
		slog.Info("record", "namespace", k, "flags", fmt.Sprintf("0b%064b", v))
	}

	// TODO: send on the wire to the listener
}
