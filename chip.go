package chip

import (
	"fmt"
	"sync/atomic"

	"github.com/google/uuid"
)

type Chip struct {
	ID     uuid.UUID
	Record atomic.Uint64
}

func NewChip(id uuid.UUID) *Chip {
	return &Chip{
		ID:     id,
		Record: atomic.Uint64{},
	}
}

func (c *Chip) Get() uint64 {
	return c.Record.Load()
}

func (c *Chip) Set(flags uint64) {
	c.Record.Or(flags)
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

	c.Record.Or(flag)

	return nil
}

func (c *Chip) MustSetPositions(position ...uint8) {
	if err := c.SetPositions(position...); err != nil {
		panic(err)
	}
}

func (c *Chip) Check(flags uint64) bool {
	return c.Record.Load()&flags == flags
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
  return c.Record.Swap(0)
}
