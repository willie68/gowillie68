package extmath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsInt(t *testing.T) {
	tsts := []struct {
		value    any
		expected any
	}{
		{
			value:    int(1234),
			expected: int(1234),
		},
		{
			value:    int(-1234),
			expected: int(1234),
		},
		{
			value:    int8(124),
			expected: int8(124),
		},
		{
			value:    int8(-124),
			expected: int8(124),
		},
		{
			value:    int16(1234),
			expected: int16(1234),
		},
		{
			value:    int16(-1234),
			expected: int16(1234),
		},
		{
			value:    int32(1234),
			expected: int32(1234),
		},
		{
			value:    int32(-1234),
			expected: int32(1234),
		},
		{
			value:    int64(1234),
			expected: int64(1234),
		},
		{
			value:    int64(-1234),
			expected: int64(1234),
		},
	}

	for _, tst := range tsts {
		var dv any
		switch v := tst.value.(type) {
		case int:
			dv = Abs(v)
		case int8:
			dv = Abs(v)
		case int16:
			dv = Abs(v)
		case int32:
			dv = Abs(v)
		case int64:
			dv = Abs(v)
		default:
			t.Fatalf("unsupported type: %T", v)
		}
		assert.Equal(t, tst.expected, dv)
	}
}
