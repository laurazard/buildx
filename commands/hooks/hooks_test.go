package hooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHintToPrint(t *testing.T) {
	testCases := []struct {
		flags       map[string]string
		numCPU      int
		hint        string
		shouldPrint bool
	}{
		{
			flags:       map[string]string{},
			numCPU:      8,
			shouldPrint: false,
		},
		{
			flags: map[string]string{
				"platform": "linux/arm64",
			},
			numCPU:      8,
			shouldPrint: true,
			hint:        multiPlatformHint,
		},
		{
			flags:       map[string]string{},
			numCPU:      2,
			shouldPrint: true,
			hint:        reduceTimeHint,
		},
		{
			flags: map[string]string{
				"platform": "linux/amd64",
			},
			numCPU:      2,
			shouldPrint: true,
			hint:        multiPlatformHint,
		},
	}

	for _, tc := range testCases {
		hint, shouldPrint := getHint("buildx", tc.flags, tc.numCPU)
		assert.Equal(t, shouldPrint, tc.shouldPrint)
		assert.Equal(t, hint, tc.hint)
	}
}
