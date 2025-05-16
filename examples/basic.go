package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/mats852/chip"
	"github.com/mats852/chip/pkg/export"
)

var wellKnownUUID = uuid.MustParse("8badf00d-cafe-beef-dead-baaaaaaaaaad")

type ShimSender struct {
	receiver *export.Receiver
}

func (s *ShimSender) Send(ctx context.Context, data []byte) error {
	return s.receiver.Handle(data)
}

type ShimReceiverRepository struct{}

func (s *ShimReceiverRepository) Store(ctx context.Context, timestamp time.Time, chip *chip.Chip) error {
	slog.Info("handling chip", "id", chip.ID, "flags", fmt.Sprintf("0b%064b", chip.Get()), "timestamp", timestamp.Format(time.RFC3339))
	return nil
}

func main() {
	chip := chip.NewChip(wellKnownUUID)

	receiver := export.NewReceiver(&ShimReceiverRepository{})

	shimSender := &ShimSender{receiver: receiver}

	exportr, err := export.NewExporter(shimSender, chip)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			pos, _ := rand.Int(rand.Reader, big.NewInt(64))
			chip.SetPositions(uint8(pos.Int64()))

			slog.Info("set", "position", pos.Int64())

			time.Sleep(250 * time.Millisecond)
		}
	}()

	if err := exportr.Serve(context.TODO()); err != nil {
		slog.Error("failed to serve", "error", err)
	}
}
