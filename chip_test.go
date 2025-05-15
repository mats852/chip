package chip

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	Flag_0001 = 0b0001
	Flag_0010 = 0b0010
	Flag_0100 = 0b0100
	Flag_1000 = 0b1000
	Flag_1x63 = 1 << 63
	Flag_1001 = 0b1001
)

func Test_Table_Set_Check(t *testing.T) {
	tests := map[string]struct {
		flags     uint64
		check     uint64
		assertion assert.BoolAssertionFunc
	}{
		"simple comparison": {
			flags:     Flag_0001,
			check:     Flag_0001,
			assertion: assert.True,
		},
		"more flags than check": {
			flags:     Flag_0001 | Flag_0010,
			check:     Flag_0001,
			assertion: assert.True,
		},
		"check is not in single flag": {
			flags:     Flag_0010,
			check:     Flag_0001,
			assertion: assert.False,
		},
		"check is not in multiple flags": {
			flags:     Flag_0010 | Flag_0100 | Flag_1000,
			check:     Flag_0001,
			assertion: assert.False,
		},
		"multiple value flag contains multiple checks": {
			flags:     Flag_1001,
			check:     Flag_0001 | Flag_1000,
			assertion: assert.True,
		},
		"believe it or not, 0 contains 0": {
			flags:     0,
			check:     0,
			assertion: assert.True,
		},
		"0 is in 0b0001 (many more zeroes after 1)": {
			flags:     Flag_0001,
			check:     0,
			assertion: assert.True,
		},
		"non-0 flag is not in 0": {
			flags:     0,
			check:     Flag_0001,
			assertion: assert.False,
		},
		"1 is not in 0 flag": {
			flags:     0,
			check:     1,
			assertion: assert.False,
		},
	}

	chip := NewChip(uuid.Nil)

	for testmsg, tt := range tests {
		t.Run(testmsg, func(t *testing.T) {
			chip.Set(tt.flags)

			tt.assertion(t, chip.Check(tt.check))

      chip.Clear()
		})
	}
}

func Test_SetPosition(t *testing.T) {
	chip := NewChip(uuid.Nil)

	chip.SetPositions(0)

	assert.False(t, chip.Check(0b10000000))
	assert.True(t, chip.Check(0b1))

	chip.Clear()

	chip.MustSetPositions(15, 0, 1, 3)

	assert.True(t, chip.Check(0b1011))
	assert.False(t, chip.Check(0b1111))

	assert.True(t, chip.MustCheckPosition(3))
}
