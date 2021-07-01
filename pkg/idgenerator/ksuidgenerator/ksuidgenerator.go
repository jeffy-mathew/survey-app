package ksuidgenerator

import "github.com/segmentio/ksuid"

type KSUIDGenerator struct {
}

func NewKSUIDGenerator() *KSUIDGenerator {
	return &KSUIDGenerator{}
}
func (g *KSUIDGenerator) Generate() ksuid.KSUID {
	return ksuid.New()
}
