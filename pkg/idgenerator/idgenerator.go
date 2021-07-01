package idgenerator

//go:generate mockgen -source=idgenerator.go -destination=./idgenerator_mock/idgenerator_mock.go -package=idgenerator_mock

import "github.com/segmentio/ksuid"

type IDGenerator interface {
	Generate() ksuid.KSUID
}
