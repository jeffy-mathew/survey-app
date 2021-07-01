package actualtimegenerator

import "time"

type ActualTimeGenerator struct {
}

func NewActualTimeGenerator() *ActualTimeGenerator {
	return &ActualTimeGenerator{}
}

func (a *ActualTimeGenerator) Now() time.Time {
	return time.Now()
}
