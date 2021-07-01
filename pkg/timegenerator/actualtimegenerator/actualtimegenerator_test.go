package actualtimegenerator

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestActualTimeGenerator_Now(t *testing.T) {
	gen := NewActualTimeGenerator()
	now := time.Now()
	generatedTime := gen.Now()
	assert.True(t, generatedTime.After(now))
}
