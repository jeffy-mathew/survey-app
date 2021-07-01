package timegenerator

//go:generate mockgen -source=timegenerator.go -destination=./timegenerator_mock/timegenerator_mock.go -package=timegenerator_mock

import "time"

//TimeGenInterface is used to introduce mock capabilities in service layer
type TimeGenInterface interface {
	Now() time.Time
}
