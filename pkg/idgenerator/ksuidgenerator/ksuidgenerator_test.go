package ksuidgenerator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKSUIDGenerator_Generate(t *testing.T) {
	gen := NewKSUIDGenerator()
	id := gen.Generate()
	assert.False(t, id.IsNil())
}
