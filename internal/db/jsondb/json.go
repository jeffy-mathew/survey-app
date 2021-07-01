package jsondb

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

type JsonDB struct {
	file *os.File
}

func NewJsonDB(fileName string) (*JsonDB, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &JsonDB{
		file: file,
	}, nil
}
func (j *JsonDB) Load(target interface{}) error {
	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(j.file)
	if err != nil {
		return err
	}
	if n == 0 {
		return nil
	}
	err = json.Unmarshal(buf.Bytes(), target)
	return err
}

func (j *JsonDB) Dump(contents interface{}) error {
	defer j.file.Close()
	err := j.file.Truncate(0)
	if err != nil {
		log.Println("error while truncating file", err)
	}
	contentsJSON, err := json.Marshal(&contents)
	if err != nil {
		return err
	}
	_, err = j.file.Write(contentsJSON)
	return err
}
