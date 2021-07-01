package timegenerator

//go:generate mockgen -source=timegenerator.go -destination=./timegenerator_mock/timegenerator_mock.go -package=timegenerator_mock

import "time"

type TimeGenInterface interface {
	Now() time.Time
}
