package main

import (
	"context"
	"crypto/rand"
	"log/slog"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/mats852/chip"
	"github.com/mats852/chip/pkg/export"
)

var wellKnownUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

func main() {
	chip := chip.NewChip(wellKnownUUID)

	exportr := export.NewExporter(chip)

	go func() {
		for {
			pos, _ := rand.Int(rand.Reader, big.NewInt(64))

			chip.SetPositions(uint8(pos.Int64()))

			slog.Info("Set position", "position", pos.Int64())

			time.Sleep(250 * time.Millisecond)
		}
	}()

	if err := exportr.Serve(context.TODO()); err != nil {
		slog.Error("failed to serve", "error", err)
	}
}
