package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/mats852/chip"
	"github.com/mats852/chip/exporter"
	"github.com/mats852/chip/receiver"
)

var wellKnownUUID = uuid.MustParse("8badf00d-cafe-beef-dead-baaaaaaaaaad")

type ShimSender struct {
	receiver *receiver.Receiver
}

func (s *ShimSender) Send(ctx context.Context, snapshot chip.Snapshot) error {
	slog.Info("sending data", "time", snapshot.Timestamp)

	data, err := json.Marshal(snapshot)
	if err != nil {
		panic(err)
	}

	return s.receiver.Handle(data)
}

type ShimReceiverRepository struct{}

func (s *ShimReceiverRepository) Store(ctx context.Context, chipSnaphot chip.Snapshot) error {
	lgr := slog.With("timestamp", chipSnaphot.Timestamp)

	for k, v := range chipSnaphot.Chips {
		lgr.With("id", k, "flags", fmt.Sprintf("0b%064b", v)).Info("received chip")
	}

	return nil
}

func main() {
	chip := chip.NewChip(wellKnownUUID)

	receiver := receiver.NewReceiver(&ShimReceiverRepository{})

	shimSender := &ShimSender{receiver: receiver}

	exporterOpts := exporter.ExporterOpts{
		Interval: 5 * time.Second,
	}

	exportr, err := exporter.NewExporter(shimSender, exporterOpts)
	if err != nil {
		panic(err)
	}

	exportr.Add(chip)

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
