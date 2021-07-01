package jsondb

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

type person struct {
	Age   int    `json:"age"`
	Place string `json:"place"`
}

func TestNewJsonDB(t *testing.T) {
	t.Run("should return jsondb with opened file", func(t *testing.T) {
		fileName := "./../../../testdata/dump-0.json"
		jsonDB, err := NewJsonDB(fileName)
		assert.NoError(t, err)
		assert.NotNil(t, jsonDB.file)
		jsonDB.file.Close()
	})
}

func TestJsonDB_Load(t *testing.T) {
	t.Run("should load entries successfully", func(t *testing.T) {
		fileName := "./../../../testdata/dump-0.json"
		jsonDB, err := NewJsonDB(fileName)
		assert.NoError(t, err)
		var target map[string]person
		err = jsonDB.Load(&target)
		assert.NoError(t, err)
		expectedEntries := map[string]person{
			"jeffy": {
				Age:   25,
				Place: "India",
			},
		}
		assert.EqualValues(t, expectedEntries, target)
	})
	t.Run("should fail with invalid input", func(t *testing.T) {
		fileName := "./../../../testdata/dump-1.json"
		jsonDB, err := NewJsonDB(fileName)
		assert.NoError(t, err)
		var target map[string]person
		err = jsonDB.Load(&target)
		assert.Error(t, err)
	})
	t.Run("should load empty json without error", func(t *testing.T) {
		fileName := "./../../../testdata/dump-2.json"
		jsonDB, err := NewJsonDB(fileName)
		assert.NoError(t, err)
		var target map[string]person
		err = jsonDB.Load(&target)
		assert.NoError(t, err)
		assert.Empty(t, target)
	})
	t.Run("should load empty file without error", func(t *testing.T) {
		fileName := "./../../../testdata/dump-3.json"
		jsonDB, err := NewJsonDB(fileName)
		assert.NoError(t, err)
		var target map[string]person
		err = jsonDB.Load(&target)
		assert.NoError(t, err)
		assert.Empty(t, target)
	})
	t.Run("should return error if file is closed", func(t *testing.T) {
		fileName := "./../../../testdata/dump-0.json"
		jsonDB, err := NewJsonDB(fileName)
		assert.NoError(t, err)
		_ = jsonDB.file.Close()
		var target map[string]person
		err = jsonDB.Load(&target)
		assert.Error(t, err)
	})
}

func TestJsonDB_Dump(t *testing.T) {
	t.Run("should dump entries successfully", func(t *testing.T) {
		fileName := "./../../../testdata/dump-4.json"
		jsonDB, err := NewJsonDB(fileName)
		assert.NoError(t, err)
		entries := map[string]person{"jeffy": {Age: 25, Place: "India"}}
		err = jsonDB.Dump(entries)
		assert.NoError(t, err)
		dumpedFile, err := os.Open(fileName)
		assert.NoError(t, err)
		defer dumpedFile.Close()
		data, err := ioutil.ReadAll(dumpedFile)
		assert.NoError(t, err)
		assert.Equal(t, `{"jeffy":{"age":25,"place":"India"}}`, string(data))
	})
	t.Run("should return error when truncate file fails", func(t *testing.T) {
		fileName := "./../../../testdata/dump-4.json"
		jsonDB, err := NewJsonDB(fileName)
		assert.NoError(t, err)
		jsonDB.file.Close()
		entries := map[string]person{"jeffy": {Age: 25, Place: "India"}}
		err = jsonDB.Dump(entries)
		assert.Error(t, err)
	})
}
