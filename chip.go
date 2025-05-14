package chip

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Chip struct {
	Records sync.Map
}

type ChipSnapshot struct {
	Records map[uint8]uint64
}

func NewChip() *Chip {
	return &Chip{
		Records: sync.Map{},
	}
}

func (c *Chip) Get(namespace uint8) uint64 {
	record, ok := c.Records.Load(namespace)
	if !ok {
		return 0
	}

	return atomic.LoadUint64(record.(*uint64))
}

func (c *Chip) Set(namespace uint8, flags uint64) {
	actual, _ := c.Records.LoadOrStore(namespace, new(uint64))
	ptr := actual.(*uint64)

	atomic.OrUint64(ptr, flags)
}

func (c *Chip) SetPositions(namespace uint8, position ...uint8) error {
	for _, p := range position {
		if p > 63 {
			return fmt.Errorf("position %d is out of range", position)
		}
	}

	actual, _ := c.Records.LoadOrStore(namespace, new(uint64))
	ptr := actual.(*uint64)

	var flag uint64
	for _, p := range position {
		flag |= 1 << p
	}

	atomic.OrUint64(ptr, flag)

	return nil
}

func (c *Chip) Check(namespace uint8, flags uint64) bool {
	record, ok := c.Records.Load(namespace)
	if !ok {
		return 0&flags == flags
	}

	return atomic.LoadUint64(record.(*uint64))&flags == flags
}

func (c *Chip) CheckPosition(namespace uint8, position uint8) bool {
	if position > 63 {
		return false
	}

	record, ok := c.Records.Load(namespace)
	if !ok {
		return false
	}

	return atomic.LoadUint64(record.(*uint64))&(1<<position) != 0
}

// Export uses Records.Range so we build the snapshot but other keys can be
// set while we are iterating over the map. In this experimental phase, I don't
// think we need to block and restart fresh every time we export.
func (c *Chip) Export() ChipSnapshot {
	snapshot := ChipSnapshot{
		Records: make(map[uint8]uint64),
	}

	c.Records.Range(func(key, value any) bool {
		snapshot.Records[key.(uint8)] = atomic.LoadUint64(value.(*uint64))

		c.Records.Delete(key)

		return true
	})

	return snapshot
}
