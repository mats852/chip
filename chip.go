package chip

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type Chip struct {
	ID        uuid.UUID
	Timestamp time.Time
	Records   sync.Map
}

func NewChip(id uuid.UUID) *Chip {
	return &Chip{
		ID:        id,
		Timestamp: time.Now(),
		Records:   sync.Map{},
	}
}

func (c *Chip) Get(ns uint8) uint64 {
	record, ok := c.Records.Load(ns)
	if !ok {
		return 0
	}

	return atomic.LoadUint64(record.(*uint64))
}

func (c *Chip) Set(ns uint8, flags uint64) {
	actual, _ := c.Records.LoadOrStore(ns, new(uint64))
	ptr := actual.(*uint64)

	atomic.OrUint64(ptr, flags)
}

func (c *Chip) SetPositions(ns uint8, pos ...uint8) error {
	for _, p := range pos {
		if p > 63 {
			return fmt.Errorf("position %d is out of range", pos)
		}
	}

	actual, _ := c.Records.LoadOrStore(ns, new(uint64))
	ptr := actual.(*uint64)

	var flag uint64
	for _, p := range pos {
		flag |= 1 << p
	}

	atomic.OrUint64(ptr, flag)

	return nil
}

func (c *Chip) Check(ns uint8, flags uint64) bool {
	record, ok := c.Records.Load(ns)
	if !ok {
		return 0&flags == flags
	}

	return atomic.LoadUint64(record.(*uint64))&flags == flags
}

func (c *Chip) CheckPosition(ns uint8, pos uint8) bool {
	if pos > 63 {
		return false
	}

	record, ok := c.Records.Load(ns)
	if !ok {
		return false
	}

	return atomic.LoadUint64(record.(*uint64))&(1<<pos) != 0
}
