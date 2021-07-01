package db

//go:generate mockgen -source=db.go -destination=./db_mock/db_mock.go -package=db_mock

type DB interface {
	Load(target interface{}) error
	Dump(contents interface{}) error
}
