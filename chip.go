package chip

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type Chip struct {
	ID    uuid.UUID
	Flags atomic.Uint64
}

func NewChip(id uuid.UUID) *Chip {
	return &Chip{
		ID:    id,
		Flags: atomic.Uint64{},
	}
}

func (c *Chip) Get() uint64 {
	return c.Flags.Load()
}

func (c *Chip) Set(flags uint64) {
	c.Flags.Or(flags)
}

func (c *Chip) SetPositions(position ...uint8) error {
	for _, p := range position {
		if p > 63 {
			return fmt.Errorf("position %d is out of range", position)
		}
	}

	var flag uint64
	for _, p := range position {
		flag |= 1 << p
	}

	c.Flags.Or(flag)

	return nil
}

func (c *Chip) MustSetPositions(position ...uint8) {
	if err := c.SetPositions(position...); err != nil {
		panic(err)
	}
}

func (c *Chip) Check(flags uint64) bool {
	return c.Flags.Load()&flags == flags
}

func (c *Chip) CheckPosition(position uint8) (bool, error) {
	if position > 63 {
		return false, fmt.Errorf("position %d is out of range", position)
	}

	return c.Check(1 << position), nil
}

func (c *Chip) MustCheckPosition(position uint8) bool {
	ok, err := c.CheckPosition(position)
	if err != nil {
		panic(err)
	}

	return ok
}

func (c *Chip) Clear() uint64 {
	return c.Flags.Swap(0)
}

type Snapshot struct {
	Timestamp time.Time
	Chips     map[uuid.UUID]uint64
}

func NewSnapshot(chips []*Chip) Snapshot {
	snapshot := Snapshot{
		Timestamp: time.Now(),
		Chips:     make(map[uuid.UUID]uint64),
	}

	for _, chip := range chips {
		snapshot.Chips[chip.ID] = chip.Clear()
	}

	return snapshot
}
