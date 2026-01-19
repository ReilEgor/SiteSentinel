package domain

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIntToUUID(t *testing.T) {
	tests := []struct {
		input    uint32
		expected uuid.UUID
	}{
		{1, uuid.MustParse("00000001-0000-0000-0000-000000000000")},
		{2, uuid.MustParse("00000002-0000-0000-0000-000000000000")},
		{3, uuid.MustParse("00000003-0000-0000-0000-000000000000")},
		{4, uuid.MustParse("00000004-0000-0000-0000-000000000000")},
		{5, uuid.MustParse("00000005-0000-0000-0000-000000000000")},
		{6, uuid.MustParse("00000006-0000-0000-0000-000000000000")},
		{7, uuid.MustParse("00000007-0000-0000-0000-000000000000")},
		{8, uuid.MustParse("00000008-0000-0000-0000-000000000000")},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("IntToUUID(%d)", tt.input), func(t *testing.T) {
			result := IntToUUID(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
